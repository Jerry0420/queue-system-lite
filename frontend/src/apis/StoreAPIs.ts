import * as httpTools from './base'

const openStore = (email: string, password: string, name: string, timezone: string, queueNames: string[]): [url: string, requestParams: RequestInit] => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
        "name": name,
        "timezone": timezone,
        "queue_names": queueNames
    })
    return [
        httpTools.generateURL("/stores"),
        {
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
            body: jsonBody 
        }
    ]
}

const signInStore = (email: string, password: string): [url: string, requestParams: RequestInit] => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
    })
    return [
        httpTools.generateURL("/stores/signin"),
        {
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
            body: jsonBody 
        }
    ]
}

const refreshToken = (): [url: string, requestParams: RequestInit] => {
    return [
        httpTools.generateURL("/stores/token"),
        {
            method: httpTools.HTTPMETHOD.PUT
        }
    ]
}

const closeStore = (storeId: number, normalToken: string): [url: string, requestParams: RequestInit] => {
    const route = "/stores/".concat(storeId.toString())
    return [
        httpTools.generateURL(route),
        {
            method: httpTools.HTTPMETHOD.DELETE,
            headers: httpTools.generateAuth(normalToken)
        }
    ]
}

const forgetPassword = (email: string): [url: string, requestParams: RequestInit] => {
    const jsonBody: string = JSON.stringify({
        "email": email,
    })
    return [
        httpTools.generateURL("/stores/password/forgot"),
        {
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
            body: jsonBody
        }
    ]
}

const updatePassword = (storeId: number, passwordToken: string, password: string): [url: string, requestParams: RequestInit] => {
    const route = "/stores/".concat(storeId.toString(), "/password")
    const jsonBody: string = JSON.stringify({
        "password_token": passwordToken,
        "password": password,
    })
    return [
        httpTools.generateURL(route), 
        {
            method: httpTools.HTTPMETHOD.PATCH,
            headers: httpTools.CONTENT_TYPE_JSON,
            body: jsonBody   
        }
    ]
}

const getStoreInfoWithSSE = (storeId: number): EventSource => {
    const route = "/stores/".concat(storeId.toString(), "/sse")
    const sse = new EventSource(httpTools.generateURL(route))
    // handle sse events outside.
    // sse.onmessage = (event) => JSON.stringify(JSON.parse(event.data))
    // sse.onopen = (event) => {}
    // sse.onerror = (event) => {}
    return sse
}

const updateStoreDescription = (storeId: number, normalToken: string, description: string): [url: string, requestParams: RequestInit] => {
    const route = "/stores/".concat(storeId.toString())
    const jsonBody: string = JSON.stringify({
        "description": description,
    })
    return [
        httpTools.generateURL(route), 
        {
            method: httpTools.HTTPMETHOD.PUT,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(normalToken)},
            body: jsonBody
        }
    ]
}

export {
    openStore,
    signInStore,
    refreshToken,
    closeStore,
    forgetPassword,
    updatePassword,
    getStoreInfoWithSSE,
    updateStoreDescription
}