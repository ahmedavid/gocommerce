const flash = document.querySelector(".flash")

console.log(flash)
if(flash) {
    setTimeout(() => {
        flash.style.display = "none"
    }, 3000)
}