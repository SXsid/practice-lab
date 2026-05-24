package src.com.sudhanshu.staticEx;

public class Human {

    public String Name;
    public int Age;
    public boolean IsMarried;
    //now it's a common btw all the values
    public static int population ;
    public Human(String Name , int Age,boolean  IsMarried){
        this.Name=Name;
        this.Age=Age;
        this.IsMarried=IsMarried;
        //INFO:
        //use class erfer while working with static vairlbe as
        //they are partof class not the objects
        Human.population+=1;
    }
    @Override
    public String toString(){
        return this.Name ;
    }
    
}
