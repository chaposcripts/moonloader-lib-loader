import './style.css';
import './app.css';

import { RequestLibsData, RequestFolderSelection, InstallSelectedLibs } from '../wailsjs/go/main/App';
import { EventsEmit, EventsOn } from '../wailsjs/runtime'
var status = 'none';
var StatusLabel = {
    none: '',
    downloading: 'Скачивание...',
    extracting: 'Распаковка...',
    copying: 'Копирование...',
    deleting: 'Удаление мусора...',
}
var log = [];

EventsOn('path:selected', (path) => document.getElementById('path').value = path.length > 0 ? path : '');
EventsOn('status:set', (newStatus) => status = newStatus);
EventsOn('list:update', (status, list, errorText) => {
    document.getElementById('loader').style.display = 'none';
    if (status) {
        document.querySelectorAll('.lib').forEach((el) => el.remove());
        Object.keys(list).forEach((name) => createLibCheckbox(name));
    } else {
        document.getElementById('error-reason').textContent = errorText ?? 'Unknown Error';
        document.getElementById('load-error').style.display = 'block';
    }
});

addEventListener('DOMContentLoaded', () => {
    RequestLibsData();
    document.getElementById('path').onclick = () => RequestFolderSelection();
    document.getElementById('install').onclick = () => InstallSelectedLibs(document.getElementById('path').value, getSelectedItems());
});

function getSelectedItems() {
    var res = Array.from(document.querySelectorAll('input[type=checkbox]:checked')).map((el) => el.id);
    console.log('selected:', res);
    return res
    // document.querySelectorAll('input[type=checkbox]').forEach((el) => console.log(el.id))
}

function createLibCheckbox(libName) {
    /*
        <div class="lib">
            <input id="SAMP.lua" type="checkbox">
            <label for="SAMP.lua">SAMP.lua</label>
        </div>
    */
    const listDiv = document.getElementById('libs-list');

    const libContainer = document.createElement('div');
    libContainer.classList.add('lib');

    const checkbox = document.createElement('input');
    checkbox.id = libName;
    checkbox.type = 'checkbox';

    const label = document.createElement('label');
    label.htmlFor = libName;
    label.textContent = libName;

    libContainer.appendChild(checkbox);
    libContainer.appendChild(label);

    listDiv.appendChild(libContainer);
}