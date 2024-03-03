# Literal value

In interinGo language, there is 3 type of object can be treat as it self <b>literally</b>, which also the result
return by all possible Evaluation process of the intepreter. They are:

- Boolean: The type boolean has two values, <code>false</code> and <code>true</code></li>
- Integer: The numbers which uses 64-bit integers</li>
- Null: Hate it or wanted it, InterinGo have null, a special object representing nothing. Both <code>null</code> and <code>false</code> make a condition false; so they can be call false values

## Example 

Boolean
```
true
```
Interger
```
510
```

> Currently there is no directed way to create a <code>NULL</code> object, here is a specific example that we can
get it

Null
```
if (1 < 0) { 1 }
```
