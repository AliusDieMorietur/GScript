struct A {
  fn a() {
    print "kek";
  }

  fn t() {
    print this.b;
  }
}

print A;

let a = A();

print a;

let b = a.a;

print b;

b();

A().a();

let a = A();

a.a();

a.b = 1;

a.t();
