# txt2stickers
## Description
Instagram made a little mistake when deploying their new generative AI sticker feature. They forgot to
implement authentication on the endpoint responsible for the generation. This program, when compiled, takes
advantage of this and allows you to generate AI stickers right from your command-line! 

Note: This program previews the stickers in the terminal after generation using Sixel. If your terminal does
not support Sixel images, the output may be a little strange. Saving the images _should_ still work.