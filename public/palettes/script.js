const charButtons = document.querySelectorAll('.char_sel_button');
const button_display = document.getElementById('button_combo');
const palette_name = document.getElementById('palette_name');
const char_portrait = document.getElementById('char_portrait');
const content = document.getElementById('cont');
const nav_buttons = document.querySelectorAll('.nav_button');
const retain_toggle = document.getElementById('retain_toggle');
const checkbox_container = document.getElementById('check_container');
const menu_toggle = document.getElementById('menu_toggle');
const menu_container = document.getElementById('menu_container');
const palette_dropdown = document.getElementById('palette_dropdown');
const splash_art = document.getElementById('splash_container');
const backgrounds = document.querySelectorAll('.bg-image');
const name_bg = document.getElementById('name_bg');
const reset_button = document.getElementById('reset_button');
const about_container = document.getElementById('about_container');
const about_content = document.getElementById('about_content');
about_button = document.getElementById('about_button');
retain_toggle.checked = false;
let firstLoad = true;
let retain_toggle_checked = retain_toggle.checked;
let palettes = {};
let selectedPalette = 1;
let activeButton;
let character = "";
const adjustedCharacters = new Set(['oleander', 'stronghoof', 'texas']);
const wideCharacters = new Set(['oleander', 'texas']);
const numberMap = new Map([
    [1, "A"],
    [2, "B"],
    [3, "C"],
    [4, "D"],
    [5, "TU"],
    [6, "VC"],
    [37, "None"],
    [39, "None"],
]);

main()

retain_toggle.addEventListener('click', () => {
    retain_toggle_checked = retain_toggle.checked;
    // console.log(`Retain: ${retain_toggle_checked}`);
    checkbox_container.classList.toggle('enabled', retain_toggle_checked);
});

menu_toggle.addEventListener('click', () => {
    menu_container.classList.toggle('hidden');
});

palette_name.addEventListener('click', () => {
    palette_dropdown.classList.toggle('hidden');
    splash_art.classList.toggle('hidden');
});

reset_button.addEventListener('click', () => {
    reset()
});

charButtons.forEach(button => {
    button.addEventListener('click', async () => {
        if (button.classList.contains('active')) return;
        character = button.id.replace('_button','');
        button.classList.add('active');
        unhide();
        makePaletteList();
        if (activeButton !== null && activeButton !== undefined) activeButton.classList.remove('active');
        activeButton = button;
        if (retain_toggle_checked) {
            if (selectedPalette >= Object.keys(palettes[character]).length) {
                selectedPalette = Object.keys(palettes[character]).length;
            }
        }
        else {
            selectedPalette = 1;
        }
        splash_art.src = `https://images.candyfloof.com/tfh-data/palettes/${character}/splash_art.png`;
        if (firstLoad) {
            firstLoad = false;
            splash_art.classList.remove('hidden');
            palette_dropdown.classList.remove('hidden');
        }
        for (let i = 0; i < backgrounds.length; i++) {
            // console.log(backgrounds[i].id);
            let character_background = backgrounds[i].id.replace('background-','');
            if (character_background === character) {
                backgrounds[i].classList.remove('hidden');
            }
            else {
                backgrounds[i].classList.add('hidden');
            }
        }
        char_portrait.className = adjustedCharacters.has(character) ? character : "" ;
        nav_buttons.forEach(button => button.classList.toggle('wide',wideCharacters.has(character)));
        palette_dropdown.scrollTop = 0;
        FillDisplays(palettes[character][selectedPalette-1].name, palettes[character][selectedPalette-1].image,selectedPalette);
    });
});

nav_buttons.forEach(button => {
    button.addEventListener('click', () => {
        const action = button.id.replace('_button','');
        switch (action) {
            case 'next':
                if (selectedPalette >= Object.keys(palettes[character]).length) {
                    selectedPalette = 1;
                }
                else {
                    selectedPalette++;
                }
                break;
            case 'prev':
                if (selectedPalette <= 1) {
                    selectedPalette = Object.keys(palettes[character]).length;
                }
                else {
                    selectedPalette--;
                }
                break;
        }
        FillDisplays(palettes[character][selectedPalette-1].name, palettes[character][selectedPalette-1].image,selectedPalette);
    })
})

async function GetPalettes() {
    try {
        const response = await axios.get(`/api/palettes/`);
        return response.data;
    }
    catch (error) {
        console.error('Error fetching palettes:', error);
        return error;
    }
}

function FillDisplays(paletteName,image,slot) {
    console.log(`${character} ${slot}: ${paletteName}`);
    if (numberMap.get(slot)) slot = numberMap.get(slot);
    button_display.innerHTML = `<img src="https://images.candyfloof.com/tfh-data/palettes/buttons/${slot}.png">`;
    char_portrait.innerHTML = `<img src="${image}">`;
    palette_name.innerHTML = paletteName.toLowerCase() === "default" ? capitalize(character) : paletteName;
}

function capitalize(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

function unhide() {
    nav_buttons.forEach(button => button.classList.remove('hidden'));
    palette_name.classList.remove('hidden');
    // checkbox_container.classList.remove('hidden');
    name_bg.classList.remove('hidden');
}

function makePaletteList(){
    palette_dropdown.innerHTML = "";
    for (let i = 1; i <= Object.keys(palettes[character]).length; i++) {
        const palette = document.createElement('div');
        palette.classList.add(`palette_option`);
        palette.id = `palette_option_${i}`;
        palette.innerText = palettes[character][i-1].name;
        palette.addEventListener('click', () => {
            const slotSelection = parseInt(palette.id.replace('palette_option_',''));
            selectedPalette = slotSelection;
            FillDisplays(palettes[character][slotSelection-1].name,palettes[character][slotSelection-1].image,slotSelection);
        });
        palette_dropdown.appendChild(palette);
    }
}

function reset() {
    if (retain_toggle_checked) {
        retain_toggle.checked = retain_toggle_checked = false;
        checkbox_container.classList.toggle('enabled', retain_toggle_checked);
    }
    if (firstLoad) return;
    selectedPalette = 1;
    nav_buttons.forEach(button => button.classList.add('hidden'));
    palette_name.classList.add('hidden');
    name_bg.classList.add('hidden');
    firstLoad = true;
    splash_art.classList.add('hidden');
    palette_dropdown.classList.add('hidden');
    palette_dropdown.innerHTML = palette_name.innerHTML = button_display.innerHTML = char_portrait.innerHTML = "";
    for (let i = 0; i < backgrounds.length; i++) {
        // console.log(backgrounds[i].id);
        let character_background = backgrounds[i].id.replace('background-','');
        if (character_background === "training") {
            backgrounds[i].classList.remove('hidden');
        }
        else {
            backgrounds[i].classList.add('hidden');
        }
    }
    charButtons.forEach(button => button.classList.remove('active'));
}

async function main() {
    palettes = await GetPalettes();
    await getAbout();
}

async function getAbout() {
    try {
        const response = await axios.get(`/api/palettes/about`);
        // const formatted = parseMarkdown(response.data);
        about_content.innerHTML = response.data;
        return response.data;
    }
    catch (error) {
        console.error('Error fetching palettes:', error);
        return error;
    }
}    

about_button.addEventListener('click', () => {
    about_container.classList.toggle('hidden');
});

about_container.addEventListener('click', (event) => {
    about_container.classList.add('hidden');
});