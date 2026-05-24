package src.com.sudhanshu.app;
import java.util.Arrays;

class Hello {
    /**
     * @param args
     */
    public static void main(String[] args) {
        Student[] students = new Student[5];
        Student Rajesh = new Student(2, "Rajesh", 89.98f);
        // untinizlie object point to null class in java
        // new key word allow the heap memeory and return the refernce to it
        System.out.println(students[0]);
        System.out.println(Rajesh.Name);
        System.out.println(Rajesh);
        System.out.println(Arrays.toString(students));
    }
}


// class is logical construct and object is physical reality(which occupy the space in heap the real
// thing we work on but
// we refer what it can have wtih the logical consruct it follow tha is class)
// new to the dynamic allocation
// INFO: funciton overslodin we can do in consturcr liek if passed only tow diff contrucor if pass
// one diff construc tn
// INFO : no idea toguht the abouve instruction but we can do that
class Student {
    int RollNo;
    // defaults can we changed
    String Name = "sid";
    float Marks;

    // bind the defult arugmet which was passed during the initiated OF THE OBJECT
    // thsi is point to the class refe we are using
    // like this is simple Rajest.RollNon this is just the object refernces
    // if we don't use this and the class protpey and varible of construtor is same the jvm is
    // confused : it's a good convesion cause sometime java get consufues sane as new this.something
    // ~ object().someting
    public Student(int rollNo, String name, float marks) {
        this.RollNo = rollNo;
        this.Name = name;
        this.Marks = marks;
    }

    void greet() {
        System.out.println("heelo my name is " + Name);
    }

}
