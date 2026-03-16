# Memorize

A simple vocab memorization app.

Usage:

```
go run main.go [filename]
```

## Usage

First configure your word list. This is a json file that is passed to `args[1]`. It should conform to the schema as described in `memorize/schema.go`

Ex:
```json
{
    "name": "Word List Name",
    "words": [
        {
            "word": "evict",
            "prompts": [
                {
                    "content": "vt. Monica's mother, who has now been __ed from her home, is staying with friends.",
                    "hint": "驱逐"
                },
                ...
            ]
        },
        ...
    ]
}
```

## UX
The app uses the SM-2 algorithm. I flash coded this for shortterm memory, so I modified the algorithm a bit for faster repetition. If you want even faster memorization tune parameters from `memorize/sm2.go`.

The app uses a simple loop to train you. It gives a cloze problem and three candidate solutions, with their vowels redacted.

```
----------------
Prompt: a. These findings add considerable weight to the claims that emotional arousal is of __ significance to relapse.
Options:
1) c__s_l
2) _bstr_ct
3) _nt_m_ly
> 
```

You need to enter the **exact same word** as in the word list. For example:

```
----------------
Prompt: vt. The President has denied the allegations, which he said were __ by his political opponents.
Options:
1) d_b_s_
2) f_br_c_t_
3) _cc_st
> 
```
The answer `fabricated` will **not** be accepted. The correct answer is `fabricate`.

A correct match will result in a `q = 5`. Partial matches (answer is a substring of word) is `q = 4`. Other answers are all given `q = 2`. A blank answer is interpreted as `q = 0`. If you are confident about your knowledge, you can also directly enter your estimated `q` score. Any `q` not in the range `{0, 1, 2, 3, 4, 5}` results in UB, but still runnable.

After your answer is given, the actual word, your score and the `hint` will be displayed. Later on a hint showing toggle will be added, but right now this allows for on-the-fly memorization of words not seen before by repeatedly submitting `q = 0`.

```
----------------
Prompt: a. These findings add considerable weight to the claims that emotional arousal is of __ significance to relapse.
Options:
1) c__s_l
2) _bstr_ct
3) _nt_m_ly
> casual
Answer: causal
Score: 2
Hint: 原因的, 因果关系的
```

Your progress will be saved to `deck.cache.json`. On default, the CLI looks for `deck.cache.json` first before falling back to `args[1]`, so delete or rename that file if you want to start with a new file.