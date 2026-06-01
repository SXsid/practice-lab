
abstract class Pet {
    abstract void sound();
}

class q6 {

    static class Dog extends Pet {
        public void sound() {
            System.out.println("Dog");
        }
    }

    static class Cat extends Pet {
        public void sound() {
            System.out.println("Cat");
        }
    }

    public static void main(String[] args) {
        Pet[] pets = { new Dog(), new Cat() };
        for (Pet p : pets) {
            p.sound();
        }
    }
}
