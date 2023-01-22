import { JSONResponse } from "./reducer"
import { validateResponseSuccess } from "./helper"

const getNormalTokenFromRefreshTokenAction = (refreshTokenActionResponse: JSONResponse | undefined | null): string => {
    if (validateResponseSuccess(refreshTokenActionResponse) === true) {
        const response: JSONResponse = refreshTokenActionResponse as JSONResponse // refreshTokenActionResponse must be JSONResponse here.
        return response["token"] as string
    }
    return ""
}

const getSessionTokenFromRefreshTokenAction = (refreshTokenActionResponse: JSONResponse | undefined | null): string => {
    if (validateResponseSuccess(refreshTokenActionResponse) === true) {
        const response: JSONResponse = refreshTokenActionResponse as JSONResponse // refreshTokenActionResponse must be JSONResponse here.
        return response["session_token"] as string
    }
    return ""
}

export {
    getNormalTokenFromRefreshTokenAction,
    getSessionTokenFromRefreshTokenAction
}