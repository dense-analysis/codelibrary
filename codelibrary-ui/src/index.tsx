import * as React from 'react'
import {createRoot} from 'react-dom/client'
import {Link, Outlet, RouterProvider, createBrowserRouter, useMatches, useNavigate, useSearchParams} from 'react-router-dom'

document.body.innerHTML = '<div id="app"></div>'

type FormDataValues = {[key: string]: string | File | undefined}
type HandleSubmitFn = (event: React.FormEvent<HTMLFormElement>) => void

/**
 * Create an onSubmit handler that captures all data from a form.
 */
const createSubmitDataHandler = (
  callback: (data: FormDataValues) => void,
): HandleSubmitFn => event => {
  event.preventDefault()

  const data: {[key: string]: string | File} = {}
  const formData = new FormData(event.currentTarget)
  formData.forEach((value, key) => { data[key] = value })

  callback(data)
}

const SearchForm = () => {
  const navigate = useNavigate()

  const handleSubmit = createSubmitDataHandler(data => {
    const q = data.q

    if (q && typeof q === 'string') {
      navigate({
        pathname: '/search',
        search: 'q=' + encodeURIComponent(q),
      })
    }
  })

  return <form onSubmit={handleSubmit}>
    <input name="q" type="text" />
    <button>Submit</button>
  </form>
}

const Breadcrumbs = () => {
  const matches = useMatches()
  const crumbs = matches
    .filter(x => x.id)
    .map((match, index, array) => {
      let label = ''
      let href = match.pathname

      if (match.id === 'root') {
        label = 'Index'
      } else if (match.id === 'search') {
        label = 'Search'
      }

      // If we're on the last item, clear the href.
      if (index === array.length - 1) {
        href = ''
      }

      return {label, href}
    })
    .filter(x => x.label)

  return <ol>
    {crumbs.map((crumb, index) => (
      <li key={index}>
        {
          crumb.href
            ? <Link to={crumb.href}>{crumb.label}</Link>
            : crumb.label
        }
      </li>
    ))}
  </ol>
}

// TODO: Generate models.ts file from Swagger docs.
interface User {
  id: string
  username: string
}

interface Language {
  id: string
  name: string
}

interface CodeSample {
  id: string
  submittedBy: User
  language: Language
  title: string
  description: string
  body: string
  created: string
  modified: string
}

interface Page<T extends object> {
  results: T[]
  count: number
}

type CodeSamplePage = Page<CodeSample>

const processResponse = async(response: Response): Promise<any> => {
  if (!response.ok) {
    const text = await response.text()

    throw new Error(`Error: ${response.status} + {text}`)
  }

  return await response.json()
}

interface CodeSampleParams {
  page: number,
  q: string,
}

const searchCodeSamples = async(params: CodeSampleParams): Promise<CodeSamplePage> => {
  const searchParams = new URLSearchParams()
  searchParams.set('page', params.page.toString())
  searchParams.set('q', params.q)

  const response = await fetch('/api/code?' + searchParams.toString())

  return processResponse(response)
}

const SearchPage = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [page, setPage] = React.useState<CodeSamplePage>({
    count: 0,
    results: [],
  })
  const query = searchParams.get('q') || ''

  React.useEffect(() => {
    const fetchParams = new URLSearchParams()
    const params: CodeSampleParams = {
      q: query,
      page: 1,
    }

    searchCodeSamples(params)
      .then(
        page => {
          setPage(page)
        },
        error => {
          console.error(error)
        }
      )
  }, [query])

  return <div>
    <div>Results: {page.count}</div>
    <ol>
      {page.results.map(sample => (
        <li>{sample.title}</li>
      ))}
    </ol>
  </div>
}

const Root = () => {
  return <div>
    <SearchForm />
    <Breadcrumbs />
    <Outlet />
  </div>
}

const router = createBrowserRouter([
  {
    path: '/',
    id: 'root',
    element: <Root />,
    errorElement: <div>things went wrong</div>,
    children: [
      {
        id: '',
        index: true,
        element: <div>index</div>,
      },
      {
        path: 'search',
        id: 'search',
        element: <SearchPage />,
      },
    ],
  },
])

const root = createRoot(document.getElementById('app')!)
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)
