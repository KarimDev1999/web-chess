package appconst

const (
	APIBasePath           = "/api"
	RouteRegister         = APIBasePath + "/register"
	RouteLogin            = APIBasePath + "/login"
	RouteTimeControls     = APIBasePath + "/time-controls"
	RouteGames            = APIBasePath + "/games"
	RouteGamesWaiting     = APIBasePath + "/games/waiting"
	RouteGamesByID        = APIBasePath + "/games/{id}"
	RouteGamesMoves       = APIBasePath + "/games/{id}/moves"
	RouteGamesJoin        = APIBasePath + "/games/{id}/join"
	RouteGamesMove        = APIBasePath + "/games/{id}/move"
	RouteGamesResign      = APIBasePath + "/games/{id}/resign"
	RouteGamesOfferDraw   = APIBasePath + "/games/{id}/draw/offer"
	RouteGamesAcceptDraw  = APIBasePath + "/games/{id}/draw/accept"
	RouteGamesDeclineDraw = APIBasePath + "/games/{id}/draw/decline"

	RouteWS = "/ws"

	SwaggerDocPath = "./docs/swagger.json"
)
