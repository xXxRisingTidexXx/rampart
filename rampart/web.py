from fastapi import FastAPI

app = FastAPI()


@app.get('/')
def _get_root():
    return [
        {
            'url': 'https://dom.ria.com/uk/realty-prodaja-kvartira-kiev-pecherskiy-lesi-ukrainki-bul-17060617.html',
            'image_urls': [
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-pecherskiy-lesi-ukrainki-bul__135692528fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-pecherskiy-lesi-ukrainki-bul__135692585fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-pecherskiy-lesi-ukrainki-bul__140511350fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-pecherskiy-lesi-ukrainki-bul__140511385fl.webp'
            ],
            'price': 88000,
            'total_area': 46,
            'living_area': 30,
            'kitchen_area': 10,
            'room_number': 2,
            'floor': 3,
            'total_floor': 5,
            'housing': 'secondary',
            'point': [30.5333129, 50.4316619],
            'ssf': 2.34552923,
            'izf': 3.10933487,
            'gzf': 1.20394283
        },
        {
            'url': 'https://dom.ria.com/uk/realty-prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia-19376478.html',
            'image_urls': [
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__138023220fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__138022923fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__138022986fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__139104075fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__139104079fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-goloseevskiy-akimovia-niya-enko-lia__139104082fl.webp'
            ],
            'price': 80000,
            'total_area': 37,
            'living_area': 25,
            'kitchen_area': 5,
            'room_number': 1,
            'floor': 6,
            'total_floor': 14,
            'housing': 'primary',
            'point': [30.4733, 50.39329],
            'ssf': 1.230348274,
            'izf': 4.329372812,
            'gzf': 0.837239403
        },
        {
            'url': 'https://dom.ria.com/uk/realty-prodaja-kvartira-kiev-osokorki-solomii-krushelnitskoy-ulitsa-19769039.html',
            'image_urls': [
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-osokorki-solomii-krushelnitskoy-ulitsa__145980630fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-osokorki-solomii-krushelnitskoy-ulitsa__145980632fl.webp',
                'https://cdn.riastatic.com/photosnew/dom/photo/prodaja-kvartira-kiev-osokorki-solomii-krushelnitskoy-ulitsa__145980633fl.webp'
            ],
            'price': 69500,
            'total_area': 37,
            'living_area': 11,
            'kitchen_area': 0,
            'room_number': 1,
            'floor': 22,
            'total_floor': 25,
            'housing': 'secondary',
            'point': [30.6494872, 50.3925324],
            'ssf': 0.921392232,
            'izf': 18.22930293,
            'gzf': 0
        }
    ]
