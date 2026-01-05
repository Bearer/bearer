fn safe_function() -> i32 {
    42
}

fn dangerous_function() -> *const i32 {
    unsafe {
        let raw_ptr = Box::into_raw(Box::new(42));
        raw_ptr
    }
}

unsafe fn fully_unsafe() -> i32 {
    let ptr = 0x1234 as *const i32;
    *ptr
}

fn main() {
    let value = safe_function();
    println!("{}", value);
    
    unsafe {
        let ptr = dangerous_function();
        println!("{}", *ptr);
    }
}

