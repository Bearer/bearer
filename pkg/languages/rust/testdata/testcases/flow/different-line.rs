fn source() -> String {
    String::from("data")
}

fn sink(data: String) {
    println!("{}", data);
}

fn main() {
    let data = source();
    sink(data);
}

