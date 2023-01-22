import * as React from 'react'
import Snackbar from '@mui/material/Snackbar'
import MuiAlert, { AlertProps } from '@mui/material/Alert'
import PropTypes from "prop-types"

const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(
  props,
  ref,
) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
})

const STATUS_TYPES = {
    INFO: 'info',
    WARNING: 'warning',
    ERROR: 'error',
    SUCCESS: 'success',
}

const StatusBar = (props) => {
    const { severity, message, setMessage } = props

    const [open, setOpen] = React.useState(false)
    const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        setOpen(false)
        setMessage('')
    }

    React.useEffect(() => {
        if (message) {
            setOpen(true)
        }
    }, [message])

    return (
        <>
            <Snackbar open={open} autoHideDuration={3000} onClose={handleClose}>
                <Alert onClose={handleClose} severity={severity} sx={{ width: '100%' }}>
                    {message}
                </Alert>
            </Snackbar>
        </>
    )
}

StatusBar.propTypes = {
    severity: PropTypes.string,
    message: PropTypes.string,
    setMessage: PropTypes.func
}

StatusBar.defaultProps = {
    severity: '',
    message: '',
    setMessage: () => {}
}

export {
    StatusBar,
    STATUS_TYPES
}
