create int list myList;
set myList = [5, 7, 8, 3, 10, 1, 2, 4, 6, 9];
println("Unsorted List:");
println(myList);
count i from 1 to len(myList) begin;
count j from i to len(myList) begin;
if (myList[j] < myList[i]) begin;
create int temp;
set temp = myList[j];
set myList[j] = myList[i];
set myList[i] = temp;
end; end; end;

println("Sorted List: ");
println(myList);