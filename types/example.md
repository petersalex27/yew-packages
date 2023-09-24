```
Array a; n: Uint =
    []: Array a; 0
    | Cons a (Array a; n): Array a; n + 1

decons: (Array a; n + 1) -> (a, (Array a; n))
decons (x::xs) = x

let arr = 
    (Cons 1) . (Cons 2) . (Cons 3) . (Cons 4) [] in -- [1, 2, 3, 4] desugared
decons arr: ?
```