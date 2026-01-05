enum Status {
    Active,
    Inactive,
    Pending,
}

enum Message {
    Quit,
    Move { x: i32, y: i32 },
    Write(String),
    ChangeColor(i32, i32, i32),
}

fn main() {
    let status = Status::Active;
    let msg = Message::Write(String::from("hello"));
}

