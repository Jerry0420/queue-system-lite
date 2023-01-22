import React, {useEffect, useState} from "react"
import { useParams, useLocation } from "react-router-dom"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { updatePassword } from "../apis/StoreAPIs"
import { StatusBar, STATUS_TYPES } from "./StatusBar"
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import Avatar from '@mui/material/Avatar'
import Typography from '@mui/material/Typography'
import TextField from '@mui/material/TextField'

const UpdatePasswordComponent = () => {
  let { storeId }: {storeId: string} = useParams()
  let location = useLocation()
  let passwordToken = ""

  if (location.search.includes("password_token")) {
    const splittedQueryStrings = location.search.split("=") 
    passwordToken = splittedQueryStrings[splittedQueryStrings.length-1]
  } else {
    passwordToken = ""
  }

  // ==================== handle all status ====================
  const [statusBarSeverity, setStatusBarSeverity] = React.useState('')
  const [statusBarMessage, setStatusBarMessage] = React.useState('')

  // 
  const [password, setPassword] = useState("")
  const [passwordAlertFlag, setPasswordAlertFlag] = useState(false)
  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setPassword(window.btoa(value))
  }

  const [confirmPassword, setConfirmPassword] = useState("")
  const [confirmPasswordAlertFlag, setConfirmPasswordAlertFlag] = useState(false)
  const handleInputConfirmPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setConfirmPassword(window.btoa(value))
  }

  const [updatePasswordAction, makeUpdatePasswordRequest] = useApiRequest(...updatePassword(parseInt(storeId), passwordToken, password))

  const doMakeUpdatePasswordRequest = () => {
    const rawPassword = window.atob(password)
    if ((8 <= rawPassword.length) && (rawPassword.length <= 15)) {
      setPasswordAlertFlag(false)
      setPassword(window.btoa(password)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
      return
    }

    const rawConfirmPassword = window.atob(confirmPassword)
    if (rawConfirmPassword === rawPassword) {
        setConfirmPasswordAlertFlag(false)
    } else {
        setConfirmPasswordAlertFlag(true)
        return
    }

    if (rawConfirmPassword && rawPassword) {
        makeUpdatePasswordRequest()
    }
  }

  useEffect(() => {
    if (updatePasswordAction.actionType === ACTION_TYPES.SUCCESS) {
      setStatusBarSeverity(STATUS_TYPES.SUCCESS)
      setStatusBarMessage("Success to update password.")
    }
    if (updatePasswordAction.actionType === ACTION_TYPES.ERROR) {
      setStatusBarSeverity(STATUS_TYPES.ERROR)
      setStatusBarMessage("Fail to update password.")
    }
  }, [updatePasswordAction.actionType])

  return (
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
            Reset Password
        </Typography>
        <Box sx={{ mt: 1 }}>
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
                name="passwordConfirm"
                label="Confirm Password"
                id="passwordConfirm"
                onChange={handleInputConfirmPassword}
                error={confirmPasswordAlertFlag}
                helperText="Type again for confirmation of the password."
            />
            <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeUpdatePasswordRequest}
                >
                Reset Password
            </Button>
        </Box>
        <StatusBar
          severity={statusBarSeverity}
          message={statusBarMessage}
          setMessage={setStatusBarMessage}
        />
    </Box>
  )
}

export {
  UpdatePasswordComponent
}