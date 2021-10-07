package cwa;
import java.util.Scanner;
public class asciiValue {

	public static void main(String[] args) {
		// TODO Auto-generated method stub
Scanner sc=new Scanner(System.in);
System.out.println("enter any character");
char a= sc.next().charAt(0);
if(a>=65 || a>=90)
	System.out.println("its capital alphabet");

else if(a>=97 && a<=122)
    System.out.println("its small alphabet");
else if(a>=47 && a<=57)
	System.out.println("its numbers");
else if (a>=48 && a<=57)
	System.out.println("its special number");
else if (a>=58 &&a<=64)
	System.out.println("its special number");
else if (a>=91 &&a<=96)
	System.out.println("its special number");
else if (a>=123 &&a<=127)
	System.out.println("its special number");
	}

}
