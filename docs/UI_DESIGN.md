# ui design & layout

## layout structure

```
  [mode]  [spinner]  [active file]

  [progress bar]
  files [n/total] · data [n/total] · [speed] · eta [duration]
```

## styling & theme

colors follow the catppuccin mocha palette:

- badges: `copy` (blue), `move` (peach), `delete` (red)
- bar gradient: `#89b4fa` (sapphire) to `#cba6f7` (mauve)
- background track: `#313244` (surface1)
- metadata & labels: `#585b70` (surface2 / muted)

## responsive width & formatting

- window resize messages (`tea.WindowSizeMsg`) adjust the progress bar width (constrained between 20 and 60 columns).
- long file paths are truncated from the left with leading ellipses (`…/foo/bar.txt`) based on terminal width.
- speed and eta are smoothed using an exponential moving average (0.15 instant + 0.85 historical).
