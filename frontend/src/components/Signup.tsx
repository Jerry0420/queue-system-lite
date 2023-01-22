import React, {useEffect, useState} from "react"
import {
  Link as RouterLink,
  useNavigate,
} from 'react-router-dom'
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import { ACTION_TYPES, useApiRequest, JSONResponse } from "../apis/reducer"
import { openStore } from "../apis/StoreAPIs"
import AddBoxIcon from '@mui/icons-material/AddBox'
import { StatusBar, STATUS_TYPES } from "./StatusBar"
import Chip from '@mui/material/Chip'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'
import Avatar from '@mui/material/Avatar'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'
import Link from '@mui/material/Link'


const SignUp = () => {
  let navigate = useNavigate()

  // ==================== handle all status ====================
  const [statusBarSeverity, setStatusBarSeverity] = React.useState('')
  const [statusBarMessage, setStatusBarMessage] = React.useState('')
  
  const timezone: string = Intl.DateTimeFormat().resolvedOptions().timeZone
  
  const [email, setEmail] = useState("")
  const [emailAlertFlag, setEmailAlertFlag] = useState(false)
  const handleInputEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setEmail(value)
  }

  const [password, setPassword] = useState("")
  const [passwordAlertFlag, setPasswordAlertFlag] = useState(false)
  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setPassword(window.btoa(value))
  }

  const [name, setName] = useState("")
  const [nameAlertFlag, setNameAlertFlag] = useState(false)
  const handleInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setName(value)
  }

  const [queueName, setQueueName] = useState("")
  const [queueNameAlertFlag, setQueueNameAlertFlag] = useState(false)
  const handleInputQueueName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setQueueName(value)
  }

  const [addQueueNameAlertFlag, setAddQueueNameAlertFlag] = useState(false)
  useEffect(() => {
    if (queueName) {
      setAddQueueNameAlertFlag(false)
    } else {
      setAddQueueNameAlertFlag(true)
    }
  }, [queueName])

  const [queueNames, setQueueNames] = useState<string[]>([])
  
  const addQueueNameToQueueNames = () => {
    const _queueNames = [...queueNames]
    _queueNames.push(queueName)
    setQueueNames(_queueNames)
    setQueueName("")
  }

  const handleDeleteQueueName = (deletedQueueName: string) => {
      var _queueNames = queueNames.filter((value, index, error): boolean => {
        return value != deletedQueueName
      })
      setQueueNames(_queueNames)
  }

  useEffect(() => {
    if (checkExistenceOfRefreshableCookie() === true) {
      const storeId = localStorage.getItem("storeId")
      if (storeId != null) {
        navigate(`/stores/${storeId}`)
      } else {
        // remove refreshable cookie for signup again
        document.cookie = "refreshable=true ; expires = Thu, 01 Jan 1970 00:00:00 GMT"
      }
    }
  }, [])

  const [openStoreAction, makeOpenStoreRequest] = useApiRequest(
    ...openStore(email, password, name, timezone, queueNames)
    )

  const doMakeOpenStoreRequest = () => {
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(email)) {
      setEmailAlertFlag(false)
    } else {
      setEmailAlertFlag(true)
      return
    }

    const rawPassword = window.atob(password)
    if ((8 <= rawPassword.length) && (rawPassword.length <= 15)) {
      setPasswordAlertFlag(false)
    } else {
      setPasswordAlertFlag(true)
      return
    }

    if (name) {
      setNameAlertFlag(false)
    } else {
      setNameAlertFlag(true)
      return
    }

    if (queueNames.length > 0) {
      setQueueNameAlertFlag(false)
    } else {
      setQueueNameAlertFlag(true)
      return
    }

    if (email && rawPassword && name && timezone && queueNames.length > 0) {
      makeOpenStoreRequest()
    }
  }

  useEffect(() => {
    if (openStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      const _jsonResponse = (openStoreAction.response as JSONResponse)
      if ((_jsonResponse["error_code"])) {
        setEmail("")
        setPassword("")
        setStatusBarSeverity(STATUS_TYPES.ERROR)
        if (_jsonResponse["error_code"] === 40901) {
          setStatusBarMessage("The store is already exist, please signin or close the store.")
        } else {
          setStatusBarMessage("Fail to find the email in account list.")
        }
      } else {
        navigate("/signin")
      }
    }
    if (openStoreAction.actionType === ACTION_TYPES.ERROR) {
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to open store.")
    }
  }, [openStoreAction.actionType])

  return (
    <Box sx={{flexGrow: 1}}>
      <Grid container component="main" sx={{ height: '100vh' }}>
        <Grid
          item
          xs={false}
          sm={false}
          md={7}
          sx={{
            backgroundImage: 'url(https://images.unsplash.com/photo-1519248200454-8f2590ed22b7?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2256&q=80)',
            backgroundRepeat: 'no-repeat',
            backgroundColor: (t) =>
              t.palette.mode === 'light' ? t.palette.grey[50] : t.palette.grey[900],
            backgroundSize: 'cover',
            backgroundPosition: 'center',
          }}
        />
        <Grid item xs={12} sm={12} md={5} component={Paper} elevation={6} square>
          <Box
            sx={{
              my: 8,
              mx: 4,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
            }}
          >
            <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }} />
            <Typography component="h1" variant="h5">
              Open Store
            </Typography>
            <Box sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                onChange={handleInputEmail}
                error={emailAlertFlag}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
                onChange={handleInputPassword}
                error={passwordAlertFlag}
                helperText="Use 8 to 15 characters with a mix of letters, numbers & symbols"
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="name"
                label="Store Name"
                type="text"
                id="name"
                autoComplete="name"
                onChange={handleInputName}
                error={nameAlertFlag}
              />
              <Grid 
                container 
                spacing={2}
                alignItems="center"
                justifyContent="flex-start"
              >
                <Grid item xs={8} sm={8}>
                  <TextField
                    fullWidth
                    required
                    margin="normal"
                    name="queueName"
                    label="Queue Name"
                    type="text"
                    id="queueName"
                    onChange={handleInputQueueName}
                    value={queueName}
                    error={queueNameAlertFlag}
                  />
                </Grid>
                <Grid item xs={4} sm={4}>
                  <Button 
                    variant="contained" 
                    startIcon={<AddBoxIcon />}
                    onClick={addQueueNameToQueueNames}
                    disabled={addQueueNameAlertFlag}
                  >
                    Add
                  </Button>
                </Grid>
              </Grid>              

              {queueNames.map((queueName: string) => (
                  <Chip 
                    sx={{ mb: 1, ml: 1, mr: 1 }}
                    label={queueName}
                    key={queueName} 
                    onDelete={() => {handleDeleteQueueName(queueName)}}
                  />
                ))}

              <Typography 
                // align="center"
                variant='subtitle2'
                sx={{mb: '-20px', mt: '20px'}}
              >
                <em>* The newly open store will be automatically closed after 24 hrs. Every Store owner can reopen their store after closing.</em>
              </Typography>
              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeOpenStoreRequest}
              >
                Open
              </Button>
              <Grid container>
                <Grid item>
                  <Link component={RouterLink} variant="body2" to="/signin">
                    {"Already have an account? Sign In"}
                  </Link>
                </Grid>
              </Grid>
            </Box>
          </Box>
        </Grid>
      </Grid>
      <StatusBar
        severity={statusBarSeverity}
        message={statusBarMessage}
        setMessage={setStatusBarMessage}
      />
    </Box>
  )
}

export {
  SignUp
}