import React, {useEffect, useContext, useState} from "react"
import { useParams, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { closeStore, getStoreInfoWithSSE, updateStoreDescription } from "../apis/StoreAPIs"
import { getNormalTokenFromRefreshTokenAction, getSessionTokenFromRefreshTokenAction } from "../apis/validator"
import { Customer, Queue, Store } from "../apis/models"
import { updateCustomer } from "../apis/CustomerAPIs"
import { AppBarWDrawer } from "./AppBarWDrawer"
import { StatusBar, STATUS_TYPES } from "./StatusBar"
import CloseIcon from '@mui/icons-material/Close'
import RefreshIcon from '@mui/icons-material/Refresh'
import ExitToAppIcon from '@mui/icons-material/ExitToApp'
import AlarmIcon from '@mui/icons-material/Alarm'
import ListItemIcon from '@mui/material/ListItemIcon'
import Divider from '@mui/material/Divider'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import Stack from '@mui/material/Stack'
import CardContent from '@mui/material/CardContent'
import CardMedia from '@mui/material/CardMedia'
import Container from '@mui/material/Container'
import Card from '@mui/material/Card'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemText from '@mui/material/ListItemText'
import TextareaAutosize from '@mui/material/TextareaAutosize'
import OutlinedInput from '@mui/material/OutlinedInput'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import Select, {SelectChangeEvent} from '@mui/material/Select'
import MenuItem from '@mui/material/MenuItem'
import DialogActions from '@mui/material/DialogActions'
import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

const StoreInfo = () => {
  let navigate = useNavigate()
  let { storeId }: {storeId: string} = useParams()

  // ==================== handle all status ====================
  const [statusBarSeverity, setStatusBarSeverity] = React.useState('')
  const [statusBarMessage, setStatusBarMessage] = React.useState('')

  // ====================== session sse ======================
  const [sessionScannedURL, setSessionScannedURL] = useState("")
  const {refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow} = useContext(RefreshTokenContext)

  useEffect(() => {
    let createSessionSSE: EventSource
    wrapCheckAuthFlow(
      () => {
        const sessionToken: string = getSessionTokenFromRefreshTokenAction(refreshTokenAction.response)
        createSessionSSE = createSessionWithSSE(sessionToken)

        createSessionSSE.onmessage = (event) => {
          setSessionScannedURL(JSON.parse(event.data)["scanned_url"])
        }
        
        createSessionSSE.onerror = (event) => {
          createSessionSSE.close()
        }
      },
      () => {
         // TODO: show error message
         // force signout
         clearCookieAndLocalstorage()
         navigate("/")
      }
    )
    return () => {
      if (createSessionSSE != null) {
        createSessionSSE.close()
      }
    }
  }, [createSessionWithSSE, refreshTokenAction.response, refreshTokenAction.exception])

  // qrcode image url
  const [qrcodeImageURL, setQrcodeImageURL] = useState("")
  useEffect(() => {
    toDataURL(sessionScannedURL, (error, url) => {
      if (url != null) {
        setQrcodeImageURL(url)
      }
    })
  }, [sessionScannedURL])

  // helper function
  const countWaitingOrProcessingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.state === 'waiting' || customer.state === 'processing')
  }
  const countWaitingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.state === 'waiting')
  }
  
  // ====================== storeinfo sse ======================
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
  }, [getStoreInfoWithSSE, setStoreInfo])

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
  const [openUpdateCustomerStateDialog, setOpenUpdateCustomerStateDialog] = React.useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null) 
  const [customerNewState, setCustomerNewState] = React.useState('')
  const handleChangeCustomerNewState = (event: SelectChangeEvent<typeof customerNewState>) => {
    setCustomerNewState(event.target.value)
  }

  const handleClickCustomerState = (customer: Customer) => {
    setOpenUpdateCustomerStateDialog(true)
    setSelectedCustomer(customer)
    setCustomerNewState(customer.state) //default state
  }

  const handleCloseCustomerStateDialog = (event: React.SyntheticEvent<unknown>, reason?: string) => {
    setOpenUpdateCustomerStateDialog(false)
    setCustomerNewState('')
    setSelectedCustomer(null)
  }

  const [updateCustomerAction, makeUpdateCustomerRequest] = useApiRequest(
    ...updateCustomer(
        selectedCustomer === null ? -1 : selectedCustomer.id,
        getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
        parseInt(storeId),
        selectedQueue === null ? -1 : selectedQueue.id,
        selectedCustomer === null ? '' : selectedCustomer.state,
        customerNewState
      )
  )
  const doMakeUpdateCustomerRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateCustomerRequest()
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  const handleUpdateCustomerNewState = () => {
    if (!customerNewState) {
      return
    }
    if ((selectedCustomer as Customer).state === customerNewState) {
      return
    }
    doMakeUpdateCustomerRequest()
  }

  useEffect(() => {
    if (updateCustomerAction.actionType === ACTION_TYPES.SUCCESS) {
      setOpenUpdateCustomerStateDialog(false)
      setCustomerNewState('')
      setStatusBarSeverity(STATUS_TYPES.SUCCESS)
      setStatusBarMessage("Success to update customer.")
    }
    if (updateCustomerAction.actionType === ACTION_TYPES.ERROR) {
      setOpenUpdateCustomerStateDialog(false)
      setCustomerNewState('')
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to update customer.")
    }
  }, [updateCustomerAction.actionType])

  // update customer state dialog
  const UpdateCustomerStateDialog = (
    ): JSX.Element => {
    return (
      <Dialog disableEscapeKeyDown open={openUpdateCustomerStateDialog} onClose={handleCloseCustomerStateDialog}>
        <DialogTitle>Update Customer State</DialogTitle>
        <DialogContent>
          <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
            <FormControl sx={{ m: 1, minWidth: 120 }}>
              <InputLabel id="dialog-select-label">State</InputLabel>
              <Select
                labelId="dialog-select-label"
                id="dialog-select"
                value={customerNewState}
                onChange={handleChangeCustomerNewState}
                input={<OutlinedInput label="State" />}
              >
                <MenuItem value="">------</MenuItem>
                <MenuItem value={'waiting'}>Waiting</MenuItem>
                <MenuItem value={'processing'}>Processing</MenuItem>
                <MenuItem value={'done'}>Done</MenuItem>
                <MenuItem value={'delete'}>Delete</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCustomerStateDialog}>Cancel</Button>
          <Button onClick={handleUpdateCustomerNewState}>Ok</Button>
        </DialogActions>
      </Dialog>
    )
  }

  // ====================== store drawer ====================== 
  // update store description
  const [openUpdateStoreDescriptionDialog, setOpenUpdateStoreDescriptionDialog] = React.useState(false)
  const [storeNewDescription, setStoreNewDescription] = React.useState('')
  const handleClickUpdateStoreDescription = () => {
    setOpenUpdateStoreDescriptionDialog(true)
    setStoreNewDescription(storeInfo.description) // default description
  }

  const handleChangeStoreNewDescription = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = event.target
    setStoreNewDescription(value)
  }

  const handleCloseUpdateDescriptionDialog = () => {
    setOpenUpdateStoreDescriptionDialog(false)
    setStoreNewDescription('')
  }

  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(
      parseInt(storeId), 
      getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
      storeNewDescription
      )
  )

  const doMakeUpdateStoreDescriptionRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateStoreDescriptionRequest()
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  const handleUpdateStoreNewDescription = () => {
    doMakeUpdateStoreDescriptionRequest()
  }

  useEffect(() => {
    if (updateStoreDescriptionAction.actionType === ACTION_TYPES.SUCCESS) {
      handleCloseUpdateDescriptionDialog()
      setStatusBarSeverity(STATUS_TYPES.SUCCESS)
      setStatusBarMessage("Success to update store description.")
    }
    if (updateStoreDescriptionAction.actionType === ACTION_TYPES.ERROR) {
      handleCloseUpdateDescriptionDialog()
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to update store description.")
    }
  }, [updateStoreDescriptionAction.actionType])

  // store drawer (update store description) 
  const UpdateStoreDescriptionDialog = (): JSX.Element => {
    return (
      <Dialog disableEscapeKeyDown open={openUpdateStoreDescriptionDialog} onClose={handleCloseUpdateDescriptionDialog}>
          <DialogTitle>Update Store</DialogTitle>
          <DialogContent>
            <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
              <FormControl sx={{ m: 1, minWidth: 120 }}>
                <TextareaAutosize
                  aria-label="empty textarea"
                  placeholder="Store Description"
                  value={storeNewDescription}
                  style={{ width: 200 }}
                  onChange={handleChangeStoreNewDescription}
                />
              </FormControl>
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseUpdateDescriptionDialog}>Cancel</Button>
            <Button onClick={handleUpdateStoreNewDescription}>Ok</Button>
          </DialogActions>
        </Dialog>
    )
  }

  // Close Store
  const [openCloseStoreDialog, setOpenCloseStoreDialog] = React.useState(false)
  const handleClickCloseStore = () => {
    setOpenCloseStoreDialog(true)
  }
  const handleCloseCloseStoreDialog = () => {
    setOpenCloseStoreDialog(false)
  }
  const [closeStoreAction, makeCloseSotreRequest] = useApiRequest(
    ...closeStore(
      parseInt(storeId),
      getNormalTokenFromRefreshTokenAction(refreshTokenAction.response) 
      )
  )

  const clearCookieAndLocalstorage = () => {
    localStorage.removeItem("storeId")
    document.cookie = "refreshable=true ; expires = Thu, 01 Jan 1970 00:00:00 GMT"
  }

  useEffect(() => {
    if (closeStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      handleCloseCloseStoreDialog()
      clearCookieAndLocalstorage()
      navigate("/")
    }
    if (closeStoreAction.actionType === ACTION_TYPES.ERROR) {
      handleCloseCloseStoreDialog()
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to close store.")
    }
  }, [closeStoreAction.actionType])

  const doMakeCloseStoreRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeCloseSotreRequest()
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }
  const handleCloseStore = () => {
    doMakeCloseStoreRequest() 
  }
  const CloseStoreDialog = (): JSX.Element => {
    return (
      <Dialog disableEscapeKeyDown open={openCloseStoreDialog} onClose={handleCloseCloseStoreDialog}>
          <DialogTitle>Close Store</DialogTitle>
          <DialogContent>
            <Typography gutterBottom align="center">
              This store will be closed and all customers' data will be sent by email.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseCloseStoreDialog}>Cancel</Button>
            <Button onClick={handleCloseStore}>Ok</Button>
          </DialogActions>
        </Dialog>
    )
  }

  // signout 
  const [openSignOutDialog, setOpenSignOutDialog] = React.useState(false)
  const handleClickSignOut = () => {
    setOpenSignOutDialog(true)
  }
  const handleCloseSignOutDialog = () => {
    setOpenSignOutDialog(false)
  }
  const handleSignOut = () => {
    handleCloseSignOutDialog()
    clearCookieAndLocalstorage()
    navigate("/") 
  }
  const SignOutDialog = (): JSX.Element => {
    return (
      <Dialog disableEscapeKeyDown open={openSignOutDialog} onClose={handleCloseSignOutDialog}>
          <DialogTitle>Sign Out?</DialogTitle>
          <DialogActions>
            <Button onClick={handleCloseSignOutDialog}>Cancel</Button>
            <Button onClick={handleSignOut}>Ok</Button>
          </DialogActions>
        </Dialog>
    )
  }

  // countdown
  const countdownTime = (): string => {
    const storeCreatedAt = new Date(storeInfo.created_at)
    const storeCloseTimeNumber = storeCreatedAt.setHours(storeCreatedAt.getHours() + 24)
    const storeCloseTime = new Date(storeCloseTimeNumber).toLocaleString()
    return storeCloseTime
  }

  const StoreDrawer = (
    <div>
      <Divider />
      <List>
        <ListItem button key={"Update Store"} onClick={handleClickUpdateStoreDescription}>
          <ListItemIcon>
            <RefreshIcon />
          </ListItemIcon>
          <ListItemText primary={"Update Store"} />
        </ListItem>
        {UpdateStoreDescriptionDialog()}

        <ListItem button key={"Close Store"} onClick={handleClickCloseStore}>
          <ListItemIcon>
            <CloseIcon />
          </ListItemIcon>
          <ListItemText primary={"Close Store"} />
        </ListItem>
        {CloseStoreDialog()}

        <ListItem button key={"Sign Out"} onClick={handleClickSignOut}>
          <ListItemIcon>
            <ExitToAppIcon />
          </ListItemIcon>
          <ListItemText primary={"Sign Out"} />
        </ListItem>
        {SignOutDialog()}

      </List>

      <Divider />
      <List>
        <ListItem key={"countdown"}>
          <ListItemIcon>
            <AlarmIcon />
          </ListItemIcon>
          <ListItemText primary={`Open till ${countdownTime()}`} />
        </ListItem>
      </List>

    </div>
  )

  return (
    <>
      <AppBarWDrawer
        storeInfo={storeInfo}
        setSelectedQueue={setSelectedQueue}
        queuesInfo={queuesInfo}
        StoreDrawer={StoreDrawer}
      >
        {selectedQueue === null && (
          <>
            <Container maxWidth="lg">
              <Container fixed>
                <Typography gutterBottom variant="h5" component="h2" align="center">
                  Please scan the QRCode to join the queue.
                </Typography>
              </Container>
              <Grid container rowSpacing={2} justifyContent="center" alignItems="center">
                <Grid item key={"all"} xs={10} sm={10} md={6}>
                  <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                    <CardMedia
                      component="img"
                      sx={{
                        display: 'block',
                        marginLeft: 'auto',
                        marginRight: 'auto',
                        width: '70%'
                      }}
                      src={qrcodeImageURL}
                      alt="qrcode image"
                    />
                    <CardContent sx={{ flexGrow: 1 }}>
                      <Typography 
                        gutterBottom 
                        style={{whiteSpace: 'pre-line'}}
                      >
                        {storeInfo.description}
                      </Typography>
                      {/* TODO: uncomment it! */}
                      {/* <a href={sessionScannedURL} target="_blank">{sessionScannedURL}</a> */}
                    </CardContent>
                  </Card>
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
                            
                            {/* Await */}
                            <TableCell align="right">{countWaitingOrProcessingCustomers(queue.customers).length}</TableCell>
                            
                            {/* next waiting */}
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
                      {waitingOrProcessingCustomersOfSelectedQueue.map((customer: Customer, index) => (
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
                              <Button onClick={() => handleClickCustomerState(customer)}>waiting</Button>
                            </TableCell>
                          )}
                          {customer.state === 'processing' && (
                            <TableCell align="right">
                              <Button sx={{color: 'red'}} onClick={() => handleClickCustomerState(customer)}>{customer.state}</Button>
                            </TableCell>
                          )}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Stack>
              {UpdateCustomerStateDialog()}
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
  StoreInfo
}