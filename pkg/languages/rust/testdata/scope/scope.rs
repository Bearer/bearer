fn scope_cursor(s: &str) -> &str { s }
fn scope_nested(s: &str) -> &str { s }
fn scope_result(s: &str) -> &str { s }

struct Request;

impl Request {
    fn body(&self) -> String {
        String::new()
    }
}

fn handler(request: Request) {
    let x = String::new();

    scope_cursor(request.body());   // detected
    scope_cursor(x + &request.body()); // not detected at cursor level

    scope_nested(request.body());     // detected
    scope_nested(x + &request.body()); // detected

    scope_result(request.body());     // detected
    scope_result(x + &request.body()); // detected
}

fn main() {
    let req: Request = Request;
    let x = String::new();

    scope_cursor(req.body());
    scope_nested(req.body());
    scope_result(req.body());
}

