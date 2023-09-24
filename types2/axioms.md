# Axioms of Type System

𝚪⊢
∈
αβψδεφγηιξκλμνοπρστθωςχυζ
ΑΒΨΔΕΦΓΗΙΞΚΛΜΝΟΠΡΣΤΘΩ΅ΧΥΖ
## Syntax
```
𝚪 = {0: Uint, ': Uint -> Uint}
1. 𝚪, x: Uint ⊢ x: Uint     [Use]
2. 𝚪 ⊢ 0: Uint              [Use]
3. 
f: Uint -> Uint
f 0 = 0
f 'x = f x

expression e = x
    | e1 e2
    | (\x → e)
    | (e1, e2)
    | e1 | e2
    | e1 of (\x -> e2) else (\y -> e3)
    | e1 = e2
    | let x = e1; e2
    | e2 where x = e1
```

## _
```
x: π in 𝚪
───────── (Var)
𝚪 ⊢ x: π
```

## →
```
 𝚪, x: A ⊢ e: B
───────────────── (Abstraction)
   (λx.e): A→B
```
```
𝚪 ⊢ e1: A→B      𝚪 ⊢ e2: A
────────────────────────────── (Application)
        𝚪 ⊢ (e1 e2): B
```

## |
```
     𝚪 ⊢ e: A
──────────────────── (Disjunction)
  𝚪 ⊢ e: (A | B)
```
```
𝚪 ⊢ e1: A    𝚪 ⊢ e2: B
────────────────────── (Expansion)
  𝚪 ⊢ (e1 | e2): A|B
```
```
𝚪 ⊢ e1: A|B  𝚪, x: A ⊢ e2: C  𝚪, y: B ⊢ e3: C
────────────────────────────────────────────── (Realization)
   𝚪 ⊢ e1 of (λx.e2) else (λy.e3): C
```

## &
```
𝚪 ⊢ e1: A    𝚪 ⊢ e2: B
────────────────────── (Construction)
   𝚪 ⊢ (e1, e2): A&B
```
```
𝚪 ⊢ (e1, e2): A&B
───────────────── (Head-Separation)
    𝚪 ⊢ e1: A
```
```
𝚪 ⊢ (e1, e2): A&B
───────────────── (Tail-Separation)
    𝚪 ⊢ e2: B
```

## forall
``` 
𝚪 ⊢ e: π    π ⊑ ρ
───────────────── (Instantiation)
    𝚪 ⊢ e: ρ
```
``` 
𝚪 ⊢ e: π    a !in free(𝚪)
───────────────────────── (Generalization)
   𝚪 ⊢ e: forall a . π
```

## where/let
```
𝚪 ⊢ e1: π    𝚪, x: π ⊢ e2: ρ
──────────────────────────── (Contextualization)
 𝚪 ⊢ (e2 where x = e1): ρ

𝚪 ⊢ e1: π    𝚪, x: π ⊢ e2: ρ
──────────────────────────── (Contextualization)
  𝚪 ⊢ (let x = e1 in e2): ρ
```

# dependent types
```
𝚪 ⊢ 𝔸: 𝚷(n: A)F    𝚪 ⊢ e: A
──────────────────────────── ()
      𝚪 ⊢ 𝔸(e) = F(e)

𝚪 ⊢ 𝔸: 𝚷(n: A)F    𝚪 ⊢ e: A
──────────────────────────── ()
      𝚪 ⊢ 𝔸(e) = F(e)
```

## ==
```
𝚪 ⊢ e1: π   π == ρ
──────────────────
    𝚪 ⊢ e1: ρ
```

## Example Proof
```
family U = {N, L0 a, L1 a, .., Ln a, ..}
A: N→U
A(n) = Ln a
A: Π(n:N)A(n)
forall x. A(n) = Ln x
𝚪 = { A: N→U, 0: N, forall a. Π(n:N)A(n) }

```