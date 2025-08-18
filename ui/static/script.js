document.addEventListener("DOMContentLoaded", function () {
  const fileInput = document.getElementById("fileInput");
  const uploadForm = document.getElementById("uploadForm");
  const fileList = document.getElementById("fileList");

  async function fetchFiles() {
    const response = await fetch("/files");
    const files = await response.json();
    renderFiles(files);
  }

  function renderFiles(files) {
    fileList.innerHTML = "";
    files.forEach(file => {
      const row = document.createElement("tr");

      const nameCell = document.createElement("td");
      nameCell.textContent = file;

      const actionCell = document.createElement("td");

      const downloadBtn = document.createElement("button");
      downloadBtn.textContent = "Download";
      downloadBtn.onclick = () => {
        window.location.href = `/download/${encodeURIComponent(file)}`;
      };

      const deleteBtn = document.createElement("button");
      deleteBtn.textContent = "Delete";
      deleteBtn.classList.add("delete");
      deleteBtn.onclick = async () => {
        await fetch(`/delete/${encodeURIComponent(file)}`, {
          method: "DELETE"
        });
        fetchFiles();
      };

      actionCell.appendChild(downloadBtn);
      actionCell.appendChild(deleteBtn);

      row.appendChild(nameCell);
      row.appendChild(actionCell);

      fileList.appendChild(row);
    });
  }

  uploadForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("file", fileInput.files[0]);

    await fetch("/upload", {
      method: "POST",
      body: formData
    });

    fileInput.value = "";
    fetchFiles();
  });

  // Initial fetch
  fetchFiles();
});
