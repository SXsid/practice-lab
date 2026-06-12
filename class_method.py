class Student:
    count = 0

    def __init__(self, name: str) -> None:
        self.name = name
        Student.count += 1

    def repr(self) -> str:
        return f"Hi i'm {self.name}"

    # BUG: in heratce is broken
    def countDatasefl(self):
        return Student.count

    @classmethod
    def countdataCls(cls):
        return cls.count

    @staticmethod
    def meta_data():
        return "hi i am a student class"


s1 = Student("sid")
s2 = Student("sid")
print(s1.repr())
print(s1.countDatasefl())
print(s1.countdataCls())

print(Student.countdataCls())

print(Student.meta_data())
