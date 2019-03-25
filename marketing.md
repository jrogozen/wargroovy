hi,

i posted here around a month ago (https://www.reddit.com/r/wargroove/comments/asdlyd/creating_a_site_to_display_custom_maps/) saying I was working on a map sharing website and looking for feedback.

it took me a bit longer than expected, but i think it's at a reasonable place to share

### features
- registration available, but not required (registration tags your username on maps you've posted. you can register afterwards, and if you're using the same browser your old maps will be correctly tagged to you)
- save your maps with multiple photo uploads, custom tags (and tag suggestions)
- text editor that supports markdown, and emoji shortcodes, e.g. :smile:
- ability for map owners to edit previously posted maps
- ability for map owners to delete their posted maps
- search by title, username, or tag
- filter by tags or type (skirmish, scenario, puzzle)
- exposes the most popular tags being used
- rate maps
- order by views, rating, name, date

### automated x-post on maps
i use the reddit API to pull content from /r/wargroove, /r/customgroove, and /r/wargroovecompetitive. i try to preserve reddit username and a link back to the reddit post. if people are willing to post in a consistent format, we could make this quite a bit more useful. here's an example of something that was pulled from reddit automatically https://wargroovy.com/maps/titan-s-cr-bibtrha23akg02a6i1e0

### known issues
- you cannot change your vote on a map once you've voted
- displayed updated date on map details page is wrong
- no one but me has used the site, so i'm sure there are a million bugs :)
- text editor does not work on in android chrome

### planned features
- improve text editor! (better emoji support, gif support, wargroove specific emojis, better markdown parsing)
- map comments
- web based map editor. i.e, being able to put terrain / map sprites on a blank canvas and save the image. useful for illustrating strategies or specific map interactions

### possible features
- guide/article hosting
- more advanced reddit api usage
- ability to post game codes per map (find/join other people playing that map)
- video uploads / embedded videos per map
- expose the api
- ???

here's the tech stack i used, for whoever is interested

#### api
- go
- postgres


#### web
- express
- react, redux

everything is hosted on gcloud. (prob way too expensive for what is actually being used :'()

i'd love to hear any feedback, comments, or feature suggestions! feel free to leave a message here, or make an issue on this github issue tracker https://github.com/jrogozen/wargroovy-public/issues.

i'm a frontend dev, so this was mostly a learning experience for the backend / api stuff. i know there are at least a few professional or hobby devs in here, so i'd love to answer any questions anyone has about the tech side of this as well :).

i also know there are a couple other websites that do map sharing already :). i already sunk a bunch of time on this site, so figured i might as well post it! if it's unecessary or no one ends up using it, it was still a valuable and fun learning experience! (and i'll still continue to work on it)


here's the site: wargroovy.com

thanks so much!
