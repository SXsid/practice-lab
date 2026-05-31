from typing import List, cast


class Animal:
    def __init__(self) -> None:
        pass

    def speak(self):
        print("i can speak")


class Dog(Animal):
    def fetch(self):
        print("fetchng dog")


class Cat(Animal):
    pass


class Bird(Animal):
    pass


def main():
    animals: List[Animal] = [Dog(), Cat(), Bird()]
    for animal in animals:
        animal.speak()
        if isinstance(animal, Dog):
            dog = cast(Dog, animal)
            dog.fetch()


main()
