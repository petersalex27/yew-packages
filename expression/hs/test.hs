module Main where

true = (\t f -> t)
recur = (\f c w -> ((c w) (\x -> x x) (true w)) ((c (f w)) (\x . x x) (true (f w))))

head = (\ls -> true ls)
tails = (\ls -> false ls)
-- (\f w -> ((w (\x . x x)) (true w)) ((f w) (\x . x x) (true (f w)))
