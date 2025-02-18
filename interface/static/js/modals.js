
function openEditSizeCategoryModal(id, name, csrfToken) {
  const modalEditSizeCategoryHTML = `
    <div class="fixed top-0 left-0 w-full h-screen flex items-center justify-center bg-black bg-opacity-10 z-50 dark:bg-opacity-50">
      <div class="bg-white rounded border border-neutral-300 px-2 py-5 w-[800px] h-[500px] overflow-y-auto shadow-sm relative dark:bg-neutral-915 dark:border-neutral-800">
        <div class='flex items-center justify-between'>
          <h3 class="font-medium flex items-center">
            <img src="/static/icons/exclamation.svg" alt="" class="w-7" />
            <span class='dark:text-neutral-500'>
              Editar categoria de tamanhos dos produtos
            </span>
          </h3>
    
          <button onclick="closeModal()" class='cursor-pointer rounded hover:bg-red-400 group'>
            <img src="/static/icons/add.svg" alt="" class='rotate-45 group-hover:brightness-[6]' />
          </button>
        </div>
        <h4 class='px-7 text-sm text-neutral-500 dark:text-neutral-600'>
          Edite os valores disponíveis abaixo.
        </h4>
        <div class="px-10 my-6">
          <form action="/produtos/tamanhos/categorias/editar/${id}" method="POST">
            <input type="text" name="csrf_token" value="${csrfToken}" hidden>
            <div class="flex flex-col my-10 mx-10">
              <label for="name" class='text-sm font-medium dark:text-neutral-500'>
                Nome da categoria
              </label>
              <input type="text" name="name" id="name" placeholder="e.g. Partes de cima" autocomplete="off" autofocus
              value="${name}"
              class="w-full p-2 bg-white border border-neutral-300 rounded font-normal focus:ring-2 focus:ring-blue-300 focus:outline-none placeholder:text-neutral-400 dark:bg-neutral-915 dark:border-neutral-800 dark:focus:ring-blue-400 dark:text-neutral-300 dark:placeholder-neutral-700" />
            </div>
        
            <div class="flex items-center justify-end absolute right-0 bottom-0 p-4">
              <button type="button" onclick="closeModal()" class="mx-5 my-2 dark:text-neutral-500">
                Cancelar
              </button>
              <button type="submit" class="rounded px-5 py-1.5 text-white bg-accent dark:bg-accent-dark dark:text-neutral-800">
                Salvar
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    `
  document.getElementById('modal-container').innerHTML = modalEditSizeCategoryHTML
}

function openEditSizeModal(id, name, sortOrder, parent, csrfToken) {
  const modalEditSizeCategoryHTML = `
    <div class="fixed top-0 left-0 w-full h-screen flex items-center justify-center bg-black bg-opacity-10 z-50 dark:bg-opacity-50">
      <div class="bg-white rounded border border-neutral-300 px-2 py-5 w-[800px] h-[500px] overflow-y-auto shadow-sm relative dark:bg-neutral-915 dark:border-neutral-800">
        <div class='flex items-center justify-between'>
          <h3 class="font-medium flex items-center">
            <img src="/static/icons/exclamation.svg" alt="" class="w-7" />
            <span class='dark:text-neutral-500'>
              Editar tamanhos dos produtos
            </span>
          </h3>

          <button onclick="closeModal()" class='cursor-pointer rounded hover:bg-red-400 group'>
            <img src="/static/icons/add.svg" alt="" class='rotate-45 group-hover:brightness-[6]' />
          </button>
        </div>
        <h4 class='px-7 text-sm text-neutral-500 dark:text-neutral-600'>
          Edite os valores disponíveis abaixo.
        </h4>
        <div class="mx-10 my-10">
          <form action="/produtos/tamanhos/editar/${id}" method="POST">
            <input type="text" name="csrf_token" value="${csrfToken}" hidden>
            <div class="flex items-end gap-4">
              <div class="flex flex-col">
                <label for="sizeName" class='text-sm font-medium dark:text-neutral-500'>
                  Nome
                </label>
                <input type="text" name="sizeName" id="sizeName" placeholder="e.g. PP" autocomplete="off" autofocus
                  value="${name}"
                  class="p-2 bg-white border border-neutral-300 rounded font-normal focus:ring-2 focus:ring-blue-300 focus:outline-none placeholder:text-neutral-400 dark:bg-neutral-915 dark:border-neutral-800 dark:focus:ring-blue-400 dark:text-neutral-300 dark:placeholder-neutral-700" />
                </div>
              <div class="flex flex-col">
                <label for="sortOrder">
                  <span class="text-sm font-medium dark:text-neutral-500">
                    Ordem de exibição
                  </span>
                </label>
                <input type="text" name="sortOrder" id="sortOrder" placeholder="e.g. 1" value="${sortOrder}"
                  autocomplete="off"
                  class="p-2 bg-white border border-neutral-300 rounded focus:ring-2 focus:ring-blue-300 focus:outline-none placeholder:text-neutral-400 dark:bg-neutral-915 dark:border-neutral-800 dark:focus:ring-blue-400 dark:text-neutral-300 dark:placeholder-neutral-700">
              </div>
              <div class="flex flex-col min-w-max">
                <label for="parent">
                  <span class="text-sm font-medium dark:text-neutral-500">
                    A qual categoria este tamanho pertence?
                  </span>
                </label>
                <input value="${parent}" disabled title="Não é possível editar a categoria da qual o tamanho pertence"
                  class="w-full p-2 bg-neutral-100 text-center border border-neutral-300 text-neutral-700 rounded dark:bg-neutral-915 dark:border-neutral-800 dark:text-neutral-300 hover:cursor-not-allowed" />
              </div>
            </div>
        
            <div class="flex items-center justify-end absolute right-0 bottom-0 p-4">
              <button type="button" onclick="closeModal()" class="mx-5 my-2 dark:text-neutral-500">
                Cancelar
              </button>
              <button type="submit" class="rounded px-5 py-1.5 text-white bg-accent dark:bg-accent-dark dark:text-neutral-800">
                Salvar
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  `
  document.getElementById('modal-container').innerHTML = modalEditSizeCategoryHTML
}

function openEditColorModal(id, name, csrfToken) {
  const modalEditSizeCategoryHTML = `
    <div class="fixed top-0 left-0 w-full h-screen flex items-center justify-center bg-black bg-opacity-10 z-50 dark:bg-opacity-50">
      <div class="bg-white rounded border border-neutral-300 px-2 py-5 w-[800px] h-[500px] overflow-y-auto shadow-sm relative dark:bg-neutral-915 dark:border-neutral-800">
        <div class='flex items-center justify-between'>
          <h3 class="font-medium flex items-center">
            <img src="/static/icons/exclamation.svg" alt="" class="w-7" />
            <span class='dark:text-neutral-500'>
              Editar cor
            </span>
          </h3>
    
          <button onclick="closeModal()" class='cursor-pointer rounded hover:bg-red-400 group'>
            <img src="/static/icons/add.svg" alt="" class='rotate-45 group-hover:brightness-[6]' />
          </button>
        </div>
        <h4 class='px-7 text-sm text-neutral-500 dark:text-neutral-600'>
          Edite os valores disponíveis abaixo.
        </h4>
        <div class="px-10 my-6">
          <form action="/produtos/cores/editar/${id}" method="POST">
            <input type="text" name="csrf_token" value="${csrfToken}" hidden>
            <div class="flex flex-col my-10 mx-10">
              <label for="name" class='text-sm font-medium dark:text-neutral-500'>
                Nome da cor
              </label>
              <input type="text" name="name" id="name" placeholder="e.g. Partes de cima" autocomplete="off" autofocus
              value="${name}"
              class="w-full p-2 bg-white border border-neutral-300 rounded font-normal focus:ring-2 focus:ring-blue-300 focus:outline-none placeholder:text-neutral-400 dark:bg-neutral-915 dark:border-neutral-800 dark:focus:ring-blue-400 dark:text-neutral-300 dark:placeholder-neutral-700" />
            </div>
        
            <div class="flex items-center justify-end absolute right-0 bottom-0 p-4">
              <button type="button" onclick="closeModal()" class="mx-5 my-2 dark:text-neutral-500">
                Cancelar
              </button>
              <button type="submit" class="rounded px-5 py-1.5 text-white bg-accent dark:bg-accent-dark dark:text-neutral-800">
                Salvar
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    `
  document.getElementById('modal-container').innerHTML = modalEditSizeCategoryHTML
}

function openImageModal(fileName) {
  const modalImageHTML = `
    <div class='fixed top-0 left-0 w-full h-full bg-neutral-800/20 z-50' onclick="closeModal()">
      <div class='flex flex-col justify-center h-full'>
        <div class="max-w-4xl mx-auto">
          <img class="object-contain max-h-96" src="https://res.cloudinary.com/loremcommerce/image/upload/${fileName}" alt="Imagem" loading="lazy">
        </div>
      </div>
    </div>
  `

  document.getElementById('modal-container').innerHTML = modalImageHTML
}

function openUpdatePhotoModal(id, profilePic, csrfToken) {
  const modalUpdatePhotoHTML = `
  <div class="fixed top-0 left-0 w-full h-screen flex items-center justify-center bg-black bg-opacity-10 z-50 dark:bg-opacity-50">
    <div class="bg-white rounded border border-neutral-300 px-2 py-5 w-[500px] h-[500px] overflow-y-auto shadow-sm relative dark:bg-neutral-915 dark:border-neutral-800">
      <div class="flex items-center justify-between">
        <h3 class="font-medium flex items-center">
          <img src="/static/icons/exclamation.svg" alt="" class="w-7" />
          <span class='dark:text-neutral-500'>
            Alterar foto de perfil
          </span>
        </h3>

        <button onclick="closeModal()" class='cursor-pointer rounded hover:bg-red-400 group'>
          <img src="/static/icons/add.svg" alt="" class='rotate-45 group-hover:brightness-[6]' />
        </button>
      </div>

      <form action="/usuarios/admin/alterar-foto/${id}" method="post" enctype="multipart/form-data"
        class="my-10 mx-auto w-fit flex flex-col items-center gap-4">
        <input type="text" name="csrf_token" value="${csrfToken}" hidden>
        <input type="file" name="image" id="image" accept="image/png, image/jpeg, image/jpg, image/webp"
          hidden>

        <label for="image" class="cursor-pointer">
          <div
            class="bg-neutral-200 border border-neutral-300 w-64 h-64 rounded-md flex flex-col overflow-hidden justify-center relative dark:bg-neutral-915 dark:border-neutral-800">
            <img src="https://res.cloudinary.com/loremcommerce/image/upload/${profilePic}" alt=""
              id="imagePreview" class="absolute w-full h-full rounded-md object-cover z-10 hidden" loading="lazy">
            <img src="/static/icons/gallery-add.svg" alt="Adicionar"
              class="w-fit h-fit mx-auto brightness-[3] dark:brightness-100">
          </div>
        </label>
        <button type="submit"
          class="bg-accent text-neutral-50 text-center rounded-md px-8 py-2 w-fit font-semibold dark:bg-accent-dark dark:text-neutral-800">
          Alterar foto
        </button>
      </form>
    </div>
  </div>
  `

  document.getElementById('modal-container').innerHTML = modalUpdatePhotoHTML
  document.getElementById('image').onchange = evt => {
    const [file] = document.getElementById('image').files
    const previewImage = document.getElementById('imagePreview')
    if (file) {
      previewImage.classList.remove('hidden')
      previewImage.src = URL.createObjectURL(file)
    }
  }
}

function closeModal() {
  document.getElementById("modal-container").innerHTML = '';
}