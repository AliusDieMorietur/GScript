let a = "global";
{
  fn showA() {
    print a;
  }
  showA();
  let a = "block";
  showA();
  print a;
}
