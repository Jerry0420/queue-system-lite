package httpAPI

import "fmt"

func V_1(route string) string {
	return fmt.Sprintf("/api/v1%s", route)
}