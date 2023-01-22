import { useReducer, useCallback } from "react"

const ACTION_TYPES = {
    RUNNING: 'running',
    SUCCESS: 'success',
    ERROR: 'error',
}

interface JSONResponse {
  [propName: string]: object | number | string | null | undefined
}

interface Action {
  actionType: string | null
  response?: JSONResponse | null
  exception?: Error | null
}

const initialState: Action = {
  actionType: null,
  response: null,
  exception: null
}

const doRunning = (): Action => ({ actionType: ACTION_TYPES.RUNNING })
const doSuccess = (response: JSONResponse): Action => ({ actionType: ACTION_TYPES.SUCCESS, response })
const doError = (exception: Error): Action => ({ actionType: ACTION_TYPES.ERROR, exception, response: null })
  
const reducer = (state = initialState, { actionType, response, exception }: Action) => {
    switch (actionType) {
      case ACTION_TYPES.RUNNING:
        return { ...initialState, actionType: ACTION_TYPES.RUNNING }
      case ACTION_TYPES.SUCCESS:
        return { ...state, actionType: ACTION_TYPES.SUCCESS, response }
      case ACTION_TYPES.ERROR:
        return { ...state, actionType: ACTION_TYPES.ERROR, response, exception }
      default:
        return state;
    }
}

const useApiRequest = (url: string, requestParams: RequestInit): [Action, () => Promise<JSONResponse | null | undefined>] => {
    const [action, dispatch] = useReducer(reducer, initialState)
  
    const makeRequest = useCallback(async () => {
      dispatch(doRunning())
      try {
        const response = await fetch(url, requestParams)
          .then(response => response.json())
          .then((jsonResponse: JSONResponse) => {
              return jsonResponse
          })
          .catch((error: Error) => {
              console.error(error)
              throw error  
          })
        dispatch(doSuccess(response))
        return response
      } catch (e) {
        const error = (e as Error)
        dispatch(doError(error))
        throw error
      }
    }, [url, requestParams]);
  
    return [action, makeRequest];
  }

export {
    ACTION_TYPES,
    Action,
    JSONResponse,
    doRunning,
    doSuccess,
    doError,
    initialState,
    reducer,
    useApiRequest
}