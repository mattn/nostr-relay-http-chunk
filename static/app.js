window.addEventListener('DOMContentLoaded', () => {
  window.fetch('/stream')
    .then((res) => {
      if (!res.ok) {
        console.log(res)
        return
      }
      const ul = document.querySelector('#result')
      const reader = res.body.getReader()
      const txt = new TextDecoder()
      const process = ({done, value}) => {
        if (done) return
        const lis = ul.querySelectorAll('li')
        if (lis.length > 3) lis[lis.length - 1].remove()
        txt.decode(value).split('\n').forEach((line) => {
          console.log(line)
          const content = JSON.parse(line).content
          const li = document.createElement('li')
          li.textContent = content
          ul.prepend(li)
        })
        return reader.read().then(process)
      }
      reader.read().then(process)
    })
    .catch((err)=>{
      console.log(err);
    })
}, false)
