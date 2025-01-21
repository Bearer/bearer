package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RequestMapping("/api")
@RestController
public class App {
  public static void main(String[] args) {
    SpringApplication.run(SpringApiApp.class, args);
  }

  @GetMapping("/base-greet")
  public String baseGreet() {
    return "Hello from the base application class!";
  }
}