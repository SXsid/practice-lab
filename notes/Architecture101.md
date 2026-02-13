## What are Backend Achitecture Desgins?
where we define the modularity ,interface and mainly the data flow of our system to satisfy the business requirements.
### Key principles
- **Modularity**: single-responsiblity components
- **Robust **- system should handle eroor really good(may be makiing a sahred library)
- **scalabel **-  independent scaleing (sysstem caan grow with tarffic)
- **Flexiblity **- desing to apabt changes (to move formad we should have to take multiple setp back)

### Key Challenges
- **Complexity** - We can overthink and desgin solution which unmaintanble and costly.
- **Adaptability** - insuring flexiblity of the system iss hard there are something which can chage some thing which will not (try to create interface and abastraction and plug and play model as much as you can ).
- **Security** -  data travelling and processign is ssafe is major and secure challenge(internal service should be private on vpc only accsible through gateways)
- **Tech Choices**- write stack to solve /build this system (always on which you have confidence)
- **Resourece Management**- Every thing has a cost sepnd visely ,sever,eng hours,manintaing,scalling .


--------
## When To Do Architecture ?
1. when project is in green field (Not effecting any other project)state 
>knowing when it's greenfield and when it's not is really important too.

2. when some small part of larger picture doens't fit right,so design a way that won't break big still fix the small.
3. when a complex problem need to get broken down into smaller pieces (should we go to a serive or use open source lib)
4. when we need to improve the effciency of the product ,imporivng code quality.



### 	""*Who ,why and how * ""
>alwsy askked this 3 question before desginin a syste5m

----


## How To impliment them 