package main;
import (
	"fmt"
	"flag"
	/*"os"
	"os/exec"*/
	"time"
	"math/rand"
	"sync"
	);

const CELULA_VIVA = 1;
const CELULA_MUERTA = 0;

type Celula struct{
	estado int
	numVecinos int
}

func main(){
	//definicion de flags
	num_gorutines := flag.Int("ng",0,"numero de goroutines");
	rows := flag.Int("r",0,"numero de filas");
	columns := flag.Int("c",0,"numero de columnas");
	iteraciones := flag.Int("i",0,"numero de iteraciones");
	metodo_particion := flag.Int("m",0,"metodo de particion : 0 - descomposicion por franja, 1 - descomposicion por bloque");
	semilla := flag.Int("s",0,"numero de celulas vivas iniciales por porcion(franja/bloque)");
	flag.Parse();
	
	fmt.Println("HILOS:", *num_gorutines);
	fmt.Println("FILAS:", *rows);
	fmt.Println("COLUMNAS:", *columns);
	fmt.Println("GENERACIONES:", *iteraciones);
	fmt.Println("METODO-PARTICION:", *metodo_particion);
	fmt.Println("SEMILLA:", *semilla);

	var wg, jo sync.WaitGroup
	jo.Add(int(*iteraciones))
	wg.Add(int(*num_gorutines))

	//array de canales de envio y de recivo
	//var channel_enviar01 = make(chan [][]Celula)
	//var channel_enviar02 = make(chan [][]Celula)

	//array de porciones de celulas del tablero
	var array_celulas_tablero [2]chan [][]Celula

	//variable para obtener las mitades de los nuevos estados de los 2 sub tableros
	var celulas = make([][]Celula, int(*rows));
	var tablero = make([][]int, int(*rows));

	//bordeIzquierdo := make([][]Celula, int(*rows))

	for i := range array_celulas_tablero {
		array_celulas_tablero[i] = make(chan [][]Celula)
	}
	//creacion del tablero
	for i := range tablero {
		tablero[i] = make([]int, int(*rows));
	}
	//creacion de las celulas del tablero
	for i := range celulas {
		celulas[i] = make([]Celula, int(*rows));
	}
	rand.Seed(time.Now().UnixNano());

	//semilla en la primera mitad del tablero
	for i := 0 ; i < int(*columns) ; i++{
		for j := 0 ; j < (int(*rows)/int(*num_gorutines)) ; j++{
			celulas[i][j].estado = rand.Intn(1 - 0 + 1) + 0;
		}	
	}
	//semilla en la segunda mitad del tablero
	for i := 0 ; i < int(*columns) ; i++{
		for j := int(*rows)/int(*num_gorutines)  ; j < int(*rows) ; j++{
			celulas[i][j].estado = rand.Intn(1 - 0 + 1) + 0;
		}	
	}
	for j := range celulas {
		fmt.Printf("%v\n",celulas[j]);		
	}
	fmt.Printf("\n");
	primeraMitad := calcularMapa(celulas,int(*num_gorutines),int(*rows),int(*columns),2)
	fmt.Printf("PRIMER BLOQUE \n");	
	for i := range primeraMitad {
		fmt.Printf("%v\n",primeraMitad[i]);		
	}
	ultimaMitad := calcularMapa(celulas,int(*num_gorutines),int(*rows),int(*columns),3)
	fmt.Printf("ULTIMO BLOQUE \n");	
	for i := range ultimaMitad {
		fmt.Printf("%v\n",ultimaMitad[i]);		
	}
	fmt.Printf("BORDE IZQUIERDO PRIMER BLOQUE \n");	
	bordeIzquierdo := obtenerBordeIzquierdo(primeraMitad,int(*rows))
	for i := range bordeIzquierdo {
		fmt.Printf("%v\n",bordeIzquierdo[i]);		
	}
	fmt.Printf("BORDE DERECHO ULTIMO BLOQUE \n");	
	bordeDerecho := obtenerBordeDerecho(ultimaMitad, int(*num_gorutines),int(*rows))
	for i := range bordeDerecho {
		fmt.Printf("%v\n",bordeDerecho[i]);		
	}

	agregarBordeFantasmaDerecho(primeraMitad,bordeDerecho)
	fmt.Printf("PRIMER BLOQUE CON BORDE FANTASMA DERECHO AGREGADO\n");	
	for i := range primeraMitad {
		fmt.Printf("%v\n",primeraMitad[i]);		
	}
	primeraMitad = agregarBordeFantasmaIzquierdo(primeraMitad,bordeDerecho)
	fmt.Printf("PRIMER BLOQUE CON BORDE FANTASMA IZQUIERDO AGREGADO\n");	
	for i := range primeraMitad {
		fmt.Printf("%v\n",primeraMitad[i]);		
	}
	
}

func calcularMapa(celulas [][]Celula, gorutines int, filas int, columnas int, k int) [][]Celula {

	bloque := columnas / gorutines

	resto := columnas % gorutines

	if resto != 0 {
		panic("Los bloques deben ser de igual tamaño, es decir, el modulo de la cantidad de columnas por la cantidad de rutinas debe ser igual a 0")
	}

	columnaMin := k * bloque
	columnaMax := (k + 1) * bloque

	newCelula := make([][]Celula, len(celulas))

	for i := range celulas {
		newCelula[i] = make([]Celula, (columnaMax - columnaMin))
		for j := 0; j < len(newCelula[i]); j++ {
			newCelula[i][j] = celulas[i][columnaMin+j]
		}
	}
	return newCelula
}

func obtenerBordeDerecho(Mitad [][]Celula, numgorutines int,filas int) [][]Celula{
	bordeDerecho := make([][]Celula, filas)
	for i := range bordeDerecho {
		bordeDerecho[i] = make([]Celula, 1);
	}
	for i := 0 ; i < filas ; i++{
		for j := 0 ; j < 1 ; j++{
			bordeDerecho[i][j] = Mitad[i][j+len(Mitad)/numgorutines-1]
		}	
	}
	return bordeDerecho

}

func obtenerBordeIzquierdo(Mitad [][]Celula,filas int)[][]Celula{
	bordeIzquierdo := make([][]Celula, filas)
	for i := range bordeIzquierdo {
		bordeIzquierdo[i] = make([]Celula, 1);
	}
	for i := 0 ; i < filas ; i++{
		for j := 0 ; j < 1 ; j++{
			bordeIzquierdo[i][j] = Mitad[i][j]
		}	
	}
	return bordeIzquierdo
}

func agregarBordeFantasmaDerecho(bloque [][]Celula,bordeDerecho [][]Celula){	
	for i := 0; i < len(bordeDerecho); i++ {
		for j := 0; j < len(bordeDerecho[i]); j++ {
			bloque[i] = append(bloque[i], bordeDerecho[i][j])
		}
	}
}

func agregarBordeFantasmaIzquierdo(bloque [][]Celula,bordeIzquierdo [][]Celula)[][]Celula{
	for i := range bloque {
		for j := range bloque[i] {
			bordeIzquierdo[i] = append(bordeIzquierdo[i], bloque[i][j])
		}
	}
	return bordeIzquierdo
}


func checkVecinos(Tablero [][]Celula,i int, j int)int{

	//Para controlar los 'runtime error: index out of range'
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	c := 0 //contador

	if Tablero[i-1][j-1].estado == CELULA_VIVA {
		c++
	}
	if Tablero[i-1][j].estado == CELULA_VIVA {
		c++
	}
	if Tablero[i-1][j+1].estado == CELULA_VIVA  {
		c++
	}

	if Tablero[i][j-1].estado == CELULA_VIVA  {
		c++
	}
	if Tablero[i][j+1].estado == CELULA_VIVA  {
		c++
	}

	if Tablero[i+1][j-1].estado == CELULA_VIVA  {
		c++
	}
	if Tablero[i+1][j].estado == CELULA_VIVA  {
		c++
	}
	if Tablero[i+1][j+1].estado == CELULA_VIVA  {
		c++
	}

	return c;

}