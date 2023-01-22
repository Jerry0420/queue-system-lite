import React, {useEffect, useState} from "react"
import { 
  Link as RouterLink,
  useNavigate 
} from "react-router-dom"
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { signInStore, forgetPassword } from "../apis/StoreAPIs"
import { StatusBar, STATUS_TYPES } from "./StatusBar"
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'
import Avatar from '@mui/material/Avatar'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'
import Link from '@mui/material/Link'
import DialogActions from '@mui/material/DialogActions'
import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'

const SignIn = () => {
  let navigate = useNavigate()

  // ==================== handle all status ====================
  const [statusBarSeverity, setStatusBarSeverity] = React.useState('')
  const [statusBarMessage, setStatusBarMessage] = React.useState('')

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

  const [signInStoreAction, makeSignInStoreRequest] = useApiRequest(...signInStore(email, password))

  const doMakeSignInStoreRequest = () => {
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
      setPassword(window.btoa(password)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
      return
    }

    if (email && rawPassword) {
      makeSignInStoreRequest()
    }
  }

  useEffect(() => {
    if (signInStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      const _jsonResponse = (signInStoreAction.response as JSONResponse)
      if ((_jsonResponse["error_code"])) {
        setEmail("")
        setPassword("")
        setStatusBarSeverity(STATUS_TYPES.ERROR)
        if (_jsonResponse["error_code"] === 40003) {
          setStatusBarMessage("Wrong password.")
        } else {
          setStatusBarMessage("Fail to find the email in account list.")
        }
      } else {
        const storeId: number = (_jsonResponse["id"] as number)
        localStorage.setItem("storeId", storeId.toString())
        navigate(`/stores/${storeId}`)
      }
    }
    if (signInStoreAction.actionType === ACTION_TYPES.ERROR) {
      setEmail("")
      setPassword("")
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to signin.")
    }
  }, [signInStoreAction.actionType])

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

  const [openForgetPasswordDialog, setOpenForgetPasswordDialog] = React.useState(false)
  const [forgetPasswordEmail, setForgetPasswordEmail] = useState("")
  const [forgetPasswordEmailAlertFlag, setForgetPasswordEmailAlertFlag] = useState(false)
  const handleInputForgetPasswordEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setForgetPasswordEmail(value)
  }
  const [forgetPasswordAction, makeForgetPasswordRequest] = useApiRequest(...forgetPassword(forgetPasswordEmail))
  const handleForgetPassword = () => {
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(forgetPasswordEmail)) {
      setForgetPasswordEmailAlertFlag(false)
    } else {
      setForgetPasswordEmailAlertFlag(true)
      return
    }
    makeForgetPasswordRequest()
  }

  useEffect(() => {
    if (forgetPasswordAction.actionType === ACTION_TYPES.SUCCESS) {
      setOpenForgetPasswordDialog(false)
      setStatusBarSeverity(STATUS_TYPES.SUCCESS)
      setStatusBarMessage("Success to send forget password email.")
    }
    if (forgetPasswordAction.actionType === ACTION_TYPES.ERROR) {
      setOpenForgetPasswordDialog(false)
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to send forget password email.")
    }
  }, [forgetPasswordAction.actionType])

  return (
    <Box sx={{flexGrow: 1}}>
      <Grid container direction="row-reverse" component="main" sx={{ height: '100vh' }}>
        <Grid
          item
          xs={false}
          sm={false}
          md={7}
          sx={{
            backgroundImage: 'url(https://images.unsplash.com/photo-1506774518161-b710d10e2733?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80)',
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
              Signin Store
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
                value={email}
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
                value={window.atob(password)}
                autoComplete="current-password"
                onChange={handleInputPassword}
                error={passwordAlertFlag}
              />     
              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeSignInStoreRequest}
              >
                Sign In
              </Button>
              <Grid container>
                <Grid item xs>
                  <Link variant="body2" sx={{"&:hover": {cursor: "pointer"}}} onClick={() => {setOpenForgetPasswordDialog(true)}}>
                    Forgot password?
                  </Link>
                  <Dialog disableEscapeKeyDown open={openForgetPasswordDialog} onClose={() => {setOpenForgetPasswordDialog(false)}}>
                    <DialogTitle>Forget Password</DialogTitle>
                    <DialogContent>
                      <TextField
                        autoFocus
                        margin="dense"
                        id="email"
                        label="Email Address"
                        type="email"
                        fullWidth
                        variant="standard"
                        autoComplete="email"
                        onChange={handleInputForgetPasswordEmail}
                        error={forgetPasswordEmailAlertFlag}
                      />
                      We'll send a link to the email for resetting password.
                    </DialogContent>
                    <DialogActions>
                      <Button onClick={() => {setOpenForgetPasswordDialog(false)}}>Cancel</Button>
                      <Button onClick={handleForgetPassword}>Ok</Button>
                    </DialogActions>
                  </Dialog>
                </Grid>
                <Grid item>
                  <Link component={RouterLink} variant="body2" to="/">
                    {"Don't have an account? Sign Up"}
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
  SignIn
}