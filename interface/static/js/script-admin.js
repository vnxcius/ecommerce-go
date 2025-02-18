const html = document.documentElement;
const themeButton = document.getElementById('theme-button');
const profileOptions = document.getElementById('profile-options');
const profileButton = document.getElementById('profile-button');
const productOptions = document.getElementById('product-options');
const productOptionsButton = document.getElementById('product-options-button');
const imageInp = document.getElementById('image');
const imagesInp = document.getElementById('images');

if (localStorage.getItem('theme') === 'dark') {
  html.classList.add('dark');
}

themeButton.addEventListener('click', () => {
  html.classList.toggle('dark');
  localStorage.setItem('theme', html.classList.contains('dark') ? 'dark' : 'light');
});

profileButton.addEventListener('click', () => {
  profileOptions.classList.toggle('hidden');
});

profileOptions.addEventListener('click', () => {
  profileOptions.classList.toggle('hidden');
});

productOptionsButton.addEventListener('click', () => {
  productOptions.classList.toggle('hidden');
  productOptions.classList.toggle('flex');
})

imageInp.onchange = evt => {
  const [file] = imageInp.files
  const previewImage = document.getElementById('imagePreview')
  if (file) {
    previewImage.classList.remove('hidden')
    previewImage.src = URL.createObjectURL(file)
  }
}

imagesInp.onchange = evt => {
  const previewImages = document.getElementById('imagesPreview')
  previewImages.innerHTML = ''
  for (const file of imagesInp.files) {
    const div = document.createElement('div')
    const img = document.createElement('img')
    div.classList.add("bg-neutral-200", "border", "border-neutral-300", "w-20", "h-20", "rounded-md", "flex", "flex-col", "overflow-hidden", "justify-center", "relative", "cursor-pointer", "dark:bg-neutral-915", "dark:border-neutral-800")
    img.src = URL.createObjectURL(file)
    img.classList.add("absolute", "w-full", "h-full", "rounded-md", "object-cover", "z-10")
    
    previewImages.appendChild(div)
    div.appendChild(img)
  }
}
