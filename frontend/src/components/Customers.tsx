import React, { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { createCustomers } from "../apis/CustomerAPIs"
import { CustomerForm } from "../apis/models"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { scanSession } from "../apis/SessionAPIs"
import { getStoreInfoWithSSE } from "../apis/StoreAPIs"
import { AppBarWDrawer } from "./AppBarWDrawer"
import { Customer, Queue, Store } from "../apis/models"
import { StatusBar, STATUS_TYPES } from "./StatusBar"
import AddBoxIcon from '@mui/icons-material/AddBox'
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack'
import Container from '@mui/material/Container'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'
import OutlinedInput from '@mui/material/OutlinedInput'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import Select, {SelectChangeEvent} from '@mui/material/Select'
import MenuItem from '@mui/material/MenuItem'
import Chip from '@mui/material/Chip'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

const CreateCustomers = () => {
  let { storeId , sessionId}: {storeId: string, sessionId: string} = useParams()

  // ==================== handle all status ====================
  const [statusBarSeverity, setStatusBarSeverity] = React.useState('')
  const [statusBarMessage, setStatusBarMessage] = React.useState('')
  
  // ================= store info sse =================
  const [storeInfo, setStoreInfo] = useState<Store>({})
  const [queuesInfo, setQueuesInfo] = useState<Queue[]>([])
  const [selectedQueue, setSelectedQueue] = useState<Queue | null>(null)

  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
      const _storeInfo = JSON.parse(event.data)
      const _queuesInfo = _storeInfo['queues']
      setStoreInfo(_storeInfo)
      setQueuesInfo(_queuesInfo)
        // console.log(JSON.parse(event.data))
    }
    
    getStoreInfoSSE.onerror = (event) => {
        getStoreInfoSSE.close()
    }

    return () => {
      if (getStoreInfoSSE != null) {
        getStoreInfoSSE.close()
      }
    }
  }, [getStoreInfoWithSSE])

  const [waitingOrProcessingCustomersOfSelectedQueue, setWaitingOrProcessingCustomersOfSelectedQueue] = useState<Customer[]>([])
  
  useEffect(() => {
    if (selectedQueue !== null) {
      const _selectedQueue = queuesInfo.filter((queue: Queue) => queue.id === selectedQueue.id)
      setSelectedQueue(_selectedQueue[0])
    }
  }, [queuesInfo])

  useEffect(() => {
    if (selectedQueue !== null) {
      if (selectedQueue.customers) {
        setWaitingOrProcessingCustomersOfSelectedQueue(
          countWaitingOrProcessingCustomers(selectedQueue.customers)
          )
      } else {
        setWaitingOrProcessingCustomersOfSelectedQueue([])
      }
    }
  }, [selectedQueue])

  // ====================== main content ======================
  // helper function
  const countWaitingOrProcessingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.state === 'waiting' || customer.state === 'processing')
  }
  const countWaitingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.state === 'waiting')
  }

  const [customerName, setCustomerName] = useState("")
  const [customerNameAlertFlag, setCustomerNameAlertFlag] = useState(false)
  const handleInputCustomerName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setCustomerName(value)
  }

  const [customerPhone, setCustomerPhone] = useState("")
  const [customerPhoneAlertFlag, setCustomerPhoneAlertFlag] = useState(false)
  const handleInputCustomerPhone = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setCustomerPhone(value)
  }

  const [customerQueueId, setCustomerQueueId] = useState(0)
  const [customerQueueIdAlertFlag, setCustomerQueueIdAlertFlag] = useState(false)
  const handleInputCustomerQueueId = (e: SelectChangeEvent<typeof customerQueueId>) => {
    setCustomerQueueId(e.target.value as number)
  }

  const [addCustomerAlertFlag, setAddCustomerAlertFlag] = useState(false)
  useEffect(() => {
    if (customerName && customerQueueId) {
      setAddCustomerAlertFlag(false)
    } else {
      setAddCustomerAlertFlag(true)
    }
  }, [customerName, customerQueueId])

  const [customersForm, setCustomersForm] = useState<CustomerForm[]>([])
  const addCustomerToCustomersForm = () => {
    const _customersForm = [...customersForm]
    _customersForm.push({
      name: customerName,
      phone: customerPhone,
      queue_id: customerQueueId
    })
    setCustomersForm(_customersForm)
    setCustomerName("")
    setCustomerPhone("")
    setCustomerQueueId(0)
    setCustomerNameAlertFlag(false)
    setCustomerPhoneAlertFlag(false)
    setCustomerQueueIdAlertFlag(false)
  }

  const showCustomerNamePhoneQueueName = (customerForm: CustomerForm): string => {
    const _selectedQueue = queuesInfo.filter((queue: Queue) => queue.id === customerForm.queue_id)
    let customerPhone = ''
    if (customerForm.phone) {
      customerPhone = customerForm.phone
    } else {
      customerPhone = '-'
    }
    return `${customerForm.name} / ${customerPhone} / ${_selectedQueue[0].name}`
  }

  const [addCustomersAlertFlag, setAddCustomersAlertFlag] = useState(false)
  const [customersFormAlertFlag, setCustomersFormAlertFlag] = useState('black')

  const handleDeleteCustomer = (deletedCustomerForm: CustomerForm) => {
    var _customerForms = customersForm.filter((customerForm, index, error): boolean => {
      return customerForm.name != deletedCustomerForm.name
    })
    setCustomersForm(_customerForms)
  }

  const [addCustomersAction, makeAddCustomersRequest] = useApiRequest(
    ...createCustomers(sessionId, parseInt(storeId), customersForm)
  )

  const doMakeAddCustomersRequest = () => {
    if (customersForm.length > 5) {
      setCustomersFormAlertFlag('red')
    } else if (1 <= customersForm.length && customersForm.length <= 5) {
      makeAddCustomersRequest()
        .then((response) => {
          setAddCustomersAlertFlag(true)
      })
    } else {
      setCustomerNameAlertFlag(true)
      setCustomerPhoneAlertFlag(true)
      setCustomerQueueIdAlertFlag(true)
    }
  }

  useEffect(() => {
    if (customersForm.length <= 5) {
      setCustomersFormAlertFlag('black')
    } else {
      setCustomersFormAlertFlag('red')
    }
  }, [customersForm])

  useEffect(() => {
    if (addCustomersAction.actionType === ACTION_TYPES.SUCCESS) {
      setStatusBarSeverity(STATUS_TYPES.SUCCESS)
      setStatusBarMessage("Success to add customers.")
    }
    if (addCustomersAction.actionType === ACTION_TYPES.ERROR) {
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to add customers.")
    }
  }, [addCustomersAction.actionType])

  // ================= scan session =================  
  const [scanSessionAction, makeScanSessionRequest] = useApiRequest(...scanSession(sessionId, parseInt(storeId)))
  useEffect(() => {
    makeScanSessionRequest()
  }, [])
  useEffect(() => {
    if (scanSessionAction.actionType === ACTION_TYPES.ERROR) {
        setAddCustomersAlertFlag(true)
    }

    // 40007: store_session exist but is already scanned.
    if ((scanSessionAction.response != null) && (scanSessionAction.response["error_code"]) && (scanSessionAction.response["error_code"] !== 40007)) {
        setAddCustomersAlertFlag(true)
    }
  }, [scanSessionAction.actionType])
  
  return (
    <>
      <AppBarWDrawer
        storeInfo={storeInfo}
        setSelectedQueue={setSelectedQueue}
        queuesInfo={queuesInfo}
        StoreDrawer={(<></>)}
      >
        {selectedQueue === null && (
          <>
            <Container maxWidth="md">
              <Container fixed>
                <Typography gutterBottom variant="h5" component="h2" align="center">
                  Please fill the form to create customers.
                  <Typography 
                    gutterBottom 
                    style={{whiteSpace: 'pre-line'}}
                  >
                    {storeInfo.description}
                  </Typography>
                </Typography>
              </Container>
              <Grid container rowSpacing={2} justifyContent="center" alignItems="center">
                <Grid item key={"all"} xs={10} sm={10} md={6}>
                  <Box sx={{ mt: 1 }}>
                    <TextField
                      margin="normal"
                      required
                      fullWidth
                      id="name"
                      label="Name"
                      name="name"
                      autoComplete="name"
                      value={customerName}
                      onChange={handleInputCustomerName}
                      error={customerNameAlertFlag}
                    />
                    <TextField
                      margin="normal"
                      fullWidth
                      name="phone"
                      label="Phone"
                      type="phone"
                      id="phone"
                      value={customerPhone}
                      autoComplete="tel"
                      onChange={handleInputCustomerPhone}
                      error={customerPhoneAlertFlag}
                    />
                    <Grid 
                      container 
                      spacing={2}
                      alignItems="center"
                      justifyContent="flex-start"
                    >
                      <Grid item xs={8} sm={8}>
                        <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
                          <FormControl 
                            sx={{ mt: 3, minWidth: 160 }} 
                            error={customerQueueIdAlertFlag}
                          >
                            <InputLabel id="queue-select-label">Queue</InputLabel>
                            <Select
                              labelId="queue-select-label"
                              id="queue-select"
                              value={customerQueueId}
                              onChange={handleInputCustomerQueueId}
                              input={<OutlinedInput label="Queue" />}
                            >
                              <MenuItem value={0}>------</MenuItem>
                              {queuesInfo.map((queue: Queue) => (
                                <MenuItem key={queue.id} value={queue.id}>{queue.name}</MenuItem>
                              ))}
                            </Select>
                          </FormControl>
                        </Box>
                      </Grid>
                      <Grid item xs={4} sm={4}>
                        <Button 
                          variant="contained" 
                          startIcon={<AddBoxIcon />}
                          onClick={addCustomerToCustomersForm}
                          disabled={addCustomerAlertFlag}
                        >
                          Add
                        </Button>
                      </Grid>
                    </Grid>              

                    {customersForm.map((customerForm: CustomerForm) => (
                        <Chip 
                          sx={{ mb: 1, ml: 1, mr: 1 }}
                          label={showCustomerNamePhoneQueueName(customerForm)}
                          key={customerForm.name} 
                          onDelete={() => {handleDeleteCustomer(customerForm)}}
                        />
                      ))}

                    <Typography 
                      align="center"
                      variant='subtitle2'
                      color={customersFormAlertFlag}
                    >
                      <em>* Can not add more than 5 customers at a time.</em>
                    </Typography>
                    <Button
                      fullWidth
                      variant="contained"
                      sx={{ mt: 3, mb: 2 }}
                      onClick={doMakeAddCustomersRequest}
                      disabled={addCustomersAlertFlag}
                    >
                      Add Customers
                    </Button>
                  </Box>
                </Grid>
                <Grid item key={"queues"} xs={12} sm={12} md={12}>
                  <TableContainer component={Paper}>
                    <Table sx={{ minWidth: '50vw' }} aria-label="simple table">
                      <TableHead>
                        <TableRow>
                          <TableCell>Queue Name</TableCell>
                          <TableCell align="right">Await</TableCell>
                          <TableCell align="right">Next</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {queuesInfo.map((queue: Queue) => (
                          <TableRow
                            key={queue.id}
                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                          >
                            <TableCell component="th" scope="row">
                              {queue.name}
                            </TableCell>

                            {/* await */}
                            <TableCell align="right">{countWaitingOrProcessingCustomers(queue.customers).length}</TableCell>
                            
                            {/* Next waiting */}
                            {countWaitingCustomers(queue.customers).length === 0 && (
                              <TableCell align="right"> - </TableCell>  
                            )}
                            {countWaitingCustomers(queue.customers).length !== 0 && (
                              <TableCell align="right">{countWaitingCustomers(queue.customers)[0].name}</TableCell>
                            )}
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                </Grid>
              </Grid>
            </Container>
          </>
        )}

        {selectedQueue !== null && (
          <>
            <Box sx={{ width: '100%' }}>
              <Stack 
                spacing={2}
                justifyContent="center"
                alignItems="center"
              >
                <Typography variant="h2" component="h2">{selectedQueue.name}</Typography>
                <TableContainer component={Paper}>
                  <Table sx={{ minWidth: '40vw' }} aria-label="simple table">
                    <TableHead>
                      <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell align="right">Phone</TableCell>
                        <TableCell align="right">State</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {waitingOrProcessingCustomersOfSelectedQueue.map((customer: Customer, index: number) => (
                        <TableRow
                          key={customer.id}
                          sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                          <TableCell component="th" scope="row">
                            [{index}] {customer.name}
                          </TableCell>

                          <TableCell align="right">
                            {customer.phone}
                          </TableCell>

                          {customer.state === 'waiting' && (
                            <TableCell align="right">
                              waiting
                            </TableCell>
                          )}
                          {customer.state === 'processing' && (
                            <TableCell align="right" sx={{color: 'red'}}>
                              {customer.state}
                            </TableCell>
                          )}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Stack>
            </Box>
          </>
        )}

      </AppBarWDrawer>
      <StatusBar
        severity={statusBarSeverity}
        message={statusBarMessage}
        setMessage={setStatusBarMessage}
      />
    </>
  )
}

export {
    CreateCustomers
}