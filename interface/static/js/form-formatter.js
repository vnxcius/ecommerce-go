function formatCpf() {
  const input = document.getElementById('cpf');
  let cpf = input.value.replace(/\D/g, ''); // Remove tudo que não é dígito

  if (cpf.length <= 11) {
    cpf = cpf.replace(/(\d{3})(\d)/, '$1.$2');
    cpf = cpf.replace(/(\d{3})(\d)/, '$1.$2');
    cpf = cpf.replace(/(\d{3})(\d{1,2})$/, '$1-$2');
  }
  input.value = cpf;
}

function formatPhone() {
  const input = document.getElementById('phone');
  let phone = input.value.replace(/\D/g, ''); // Remove tudo que não é dígito
  if (phone.length <= 11) {
    phone = phone.replace(/^(\d{2})(\d)/, '($1) $2');
    phone = phone.replace(/(\d{5})(\d)/, '$1-$2');
  }
  input.value = phone;
}

function formatCep() {
  const input = document.getElementById('cep');
  let cep = input.value.replace(/\D/g, ''); // Remove tudo que não é dígito
  if (cep.length <= 8) {
    cep = cep.replace(/^(\d{5})(\d)/, '$1-$2');
  }
  input.value = cep;
}

document.getElementById('form').addEventListener('submit', function (e) {
  e.preventDefault();
  // Captura os campos do formulário
  let cpf = document.getElementById('cpf')?.value;
  let phone = document.getElementById('phone')?.value;
  let cep = document.getElementById('cep')?.value;

  // Remove caracteres não numéricos
  if (cpf) cpf = cpf.replace(/\D/g, '');
  if (phone) phone = phone.replace(/\D/g, '');
  if (cep) cep = cep.replace(/\D/g, '');

  if (cep) document.getElementById('cep').value = cep;
  if (cpf) document.getElementById('cpf').value = cpf;
  if (phone) document.getElementById('phone').value = phone;
  this.submit();
})

// Definir data de nascimento padrão
document.getElementById("birthDate").defaultValue = "2024-01-01";