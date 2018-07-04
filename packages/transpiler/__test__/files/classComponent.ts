@Component({})
class ClassComponent {
  constructor(private readonly abc: string) {}
  hello(x: string, y: number) {
    console.log(x, y)
  }
}
