namespace py example

struct Response {
    1: i64 cliID
    2: string operator
    3: optional string buffer
}

struct Request {
    1: i64 cliID
    2: string operator
}

service GetLock {
    Response dolock(1:Request req),
    Response unlock(1:Request req),
    Response client_states(),
}