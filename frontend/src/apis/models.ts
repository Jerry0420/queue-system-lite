interface CustomerForm {
    name: string
    phone: string
    queue_id: number
}

interface Customer {
    created_at: string
    id: number
    name: string
    phone: string
    state: string
}

interface Queue {
    customers: Customer[]
    id: number
    name: string
}

interface Store {
    created_at: string
    description: string
    email: string
    id: number
    name: string
    queues: Queue[]
}

export {
    CustomerForm,
    Customer,
    Queue,
    Store
}