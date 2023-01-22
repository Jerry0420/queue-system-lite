import * as httpTools from './base'
import {CustomerForm} from './models'

const createCustomers = (sessionId: string, storeId: number, customers: CustomerForm[]): [url: string, requestParams: RequestInit] => {
    const jsonBody: string = JSON.stringify({
        "store_id": storeId,
        "customers": customers
    })
    return [
        httpTools.generateURL("/customers"),
        {
            method: httpTools.HTTPMETHOD.POST,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(sessionId, false)},
            body: jsonBody
        }
    ]
}

const updateCustomer = (customerId: number, normalToken: string, storeId: number, queueId: number, oldCustomerState: string, newCustomerState: string): [url: string, requestParams: RequestInit] => {
    const route = "/customers/".concat(customerId.toString())
    const jsonBody: string = JSON.stringify({
        "store_id": storeId,
        "queue_id": queueId,
        "old_customer_state": oldCustomerState,
        "new_customer_state": newCustomerState,
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
    createCustomers,
    updateCustomer
}