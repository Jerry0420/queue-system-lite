import React, {useEffect} from "react"
import { useNavigate } from 'react-router-dom'

const NoMatch = () => {
    let navigate = useNavigate()
    useEffect(() => {
        navigate("/")
      }, [])
    return (<></>)
}

export {
    NoMatch
}