struct A {
  fn a() {
    print "kek";
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

