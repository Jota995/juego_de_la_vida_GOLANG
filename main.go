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

	var wg sync.WaitGroup
	wg.Add(int(*num_gorutines))

	//array de canales de envio y de recivo
	var array_channel_enviar []chan int
	var array_channel_recivir []chan int

	//array de porciones de celulas del tablero
	var array_celulas_tablero []chan [][]Celula

	var MitadCelulas1 = make([][]Celula, (int(*rows)/int(*num_gorutines)));
	var MitadCelulas2 = make([][]Celula, (int(*rows)/int(*num_gorutines)));
	//variable para obtener las mitades de los nuevos estados de los 2 sub tableros
	var celulasmitad = make([][]Celula, int(*rows)/int(*num_gorutines));
	var celulas = make([][]Celula, int(*rows));
	var tablero = make([][]int, int(*rows));


	//creacion para 2 gorutines
	if int(*num_gorutines) == 2 {

		for i := 0 ; i < int(*num_gorutines) ; i++{
			array_channel_enviar = append(array_channel_enviar, make(chan int,int(*columns)))
		}

		for i := 0 ; i < int(*num_gorutines) ; i++{
			array_channel_enviar = append(array_channel_enviar, make(chan int,int(*columns)))
		}

		for i := 0 ; i < int(*num_gorutines) ; i++{
			array_celulas_tablero = append(array_celulas_tablero, make(chan [][]Celula))
		}

		//variable para almacenar la mitad izquierda del tablero
		for i := range MitadCelulas1 {
			MitadCelulas1[i] = make([]Celula, (int(*rows)/int(*num_gorutines))+1);
		}
		//variable para almacenar la mitad derecha del tablero
		for i := range MitadCelulas2 {
			MitadCelulas2[i] = make([]Celula, (int(*rows)/int(*num_gorutines))+1);
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

		//estados de la primera mitad del tablero
		for i := 0 ; i < int(*columns)/2 ; i++{
			for j := 0 ; j < (int(*rows)/2) ; j++{
				MitadCelulas1[i][j].estado = celulas[i][j].estado
			}	
		}

		//estados de la segunda mitad del tablero
		for i := 0 ; i < int(*columns)/2 ; i++{
			for j := 1 ; j < (int(*rows)/2)+1 ; j++{
				MitadCelulas2[i][j].estado = celulas[i][j].estado
			}	
		}

		for i := 0 ; i < int(*iteraciones) ; i++{
			go func() {
				defer wg.Done()
				//recivimos el borde izquierdo del segundo sub tablero
				for i := 0 ; i < (int(*rows)/int(*num_gorutines)) ; i++{
					for j := (int(*rows)/int(*num_gorutines))-1 ; j < (int(*rows)/int(*num_gorutines)) ; j++{
						MitadCelulas1[i][j].estado = <- array_channel_recivir[1]
					}	
				}
				//check a todos los vecinos de cada celula
				for j := 0 ; j < int(*columns)/int(*num_gorutines) ; j++{
					for k := 0 ; k < int(*rows)/int(*num_gorutines) ; k++{
						MitadCelulas1[j][k].numVecinos = checkVecinos(MitadCelulas1,j,k);
					}	
				}
				//cambiar el estado de las celulas
				for i := 0 ; i < int(*columns)/int(*num_gorutines) ; i++{
					for j := 0 ; j < int(*rows)/int(*num_gorutines) ; j++{
						if ((celulas[i][j].estado == CELULA_MUERTA && celulas[i][j].numVecinos >=3) || (celulas[i][j].estado == CELULA_VIVA && (celulas[i][j].numVecinos == 2 || celulas[i][j].numVecinos == 3 ))) {
							celulas[i][j].estado = CELULA_VIVA;
						}else{
							celulas[i][j].estado = CELULA_MUERTA;
						}
					}	
				}
		
				//pasar borde derecho de la primera mitad del tablero
				for i := 0 ; i < int(*columns) ; i++{
					for j := (int(*rows)/int(*num_gorutines))-1 ; j < (int(*rows)/int(*num_gorutines)) ; j++{
						array_channel_enviar[1] <- celulas[i][j].estado
					}	
				}
		
				//pasar el nuevo estado del primer lado del tablero
				array_celulas_tablero[0] <- MitadCelulas1
		
			}()
			
			go func() {
				defer wg.Done()
				// pasar la fila derecha de la primera parte del tablero
				for i := 0 ; i < (int(*columns)/int(*num_gorutines)) ; i++{
					for j := 0 ; j < (int(*rows)/int(*num_gorutines)) - ((int(*rows)/int(*num_gorutines))-1) ; j++{
						MitadCelulas2[i][j].estado = <- array_channel_recivir[0]
					}	
				}
		
				//check a todos los vecinos de cada celula
				for j := 0 ; j < int(*columns)/int(*num_gorutines) ; j++{
					for k := 0 ; k < int(*rows)/int(*num_gorutines) ; k++{
						MitadCelulas2[j][k].numVecinos = checkVecinos(MitadCelulas2,j,k);
					}	
				}
				//cambiar el estado de las celulas
				for i := 0 ; i < int(*columns)/int(*num_gorutines) ; i++{
					for j := 0 ; j < int(*rows)/int(*num_gorutines) ; j++{
						if ((celulas[i][j].estado == CELULA_MUERTA && celulas[i][j].numVecinos >=3) || (celulas[i][j].estado == CELULA_VIVA && (celulas[i][j].numVecinos == 2 || celulas[i][j].numVecinos == 3 ))) {
							celulas[i][j].estado = CELULA_VIVA;
						}else{
							celulas[i][j].estado = CELULA_MUERTA;
						}
					}	
				}
				//pasar borde izquierdo de la segunda mitad del tablero
				for i := 0 ; i < int(*columns) ; i++{
					for j := (int(*rows)/int(*num_gorutines)) ; j < (int(*rows)/int(*num_gorutines))+1; j++{
						array_channel_enviar[2] <- celulas[i][j].estado
					}	
				}
				//pasar el nuevo estado del segundo lado del tablero
				array_celulas_tablero[1] <- MitadCelulas2
		
			}()

			wg.Wait()

			//obtener la primera mitad del tablero
			celulasmitad = <- array_celulas_tablero[1]
			for i := 0 ; i < int(*columns) ; i++{
				for j := 0 ; j < (int(*rows)/int(*num_gorutines)) ; j++{
					celulas[i][j].estado = celulasmitad[i][j].estado;
				}	
			}

			//obtener la segunda mitad del tablero
			celulasmitad = <- array_celulas_tablero[2]
			for i := 0 ; i < int(*columns) ; i++{
				for j := int(*rows)/int(*num_gorutines) ; j < int(*rows) ; j++{
					celulas[i][j].estado = celulasmitad[i][j].estado;
				}	
			}
		}	
	}

	/*
	//imprimir el primer estado del tablero
	for j := range tablero {
		fmt.Printf("%v\n",tablero[j]);		
	}
	//pausar por 2 segundos
	time.Sleep(2 * time.Second);
	c := exec.Command("clear");
	c.Stdout = os.Stdout;
	c.Run();

	for i := 0 ; i<int(*iteraciones); i++{
		fmt.Printf("Generacion : %v\n",i);

		for j := range tablero {
			fmt.Printf("%v\n",tablero[j]);		
		}
		//check a todos los vecinos de cada celula
		for j := 0 ; j < int(*columns) ; j++{
			for k := 0 ; k < int(*rows) ; k++{
				celulas[j][k].numVecinos = checkVecinos(tablero,j,k);
			}	
		}
		//cambiar el estado de las celulas
		for i := 0 ; i < int(*columns) ; i++{
			for j := 0 ; j < int(*rows) ; j++{
				if ((celulas[i][j].estado == CELULA_MUERTA && celulas[i][j].numVecinos >=3) || (celulas[i][j].estado == CELULA_VIVA && (celulas[i][j].numVecinos == 2 || celulas[i][j].numVecinos == 3 ))) {
					celulas[i][j].estado = CELULA_VIVA;
				}else{
					celulas[i][j].estado = CELULA_MUERTA;
				}
			}	
		}

		// pasar el estado de las celulas al tablero
		for i := 0 ; i < int(*columns) ; i++{
			for j := 0 ; j < int(*rows) ; j++{
				tablero[i][j] = celulas[i][j].estado;
			}	
		}

		time.Sleep(2 * time.Second)
		c := exec.Command("clear");
		c.Stdout = os.Stdout;
		c.Run();

	}

	*/
	

	/*for i := 0 ; i < int(*columns) ; i++{
		for j := 2 ; j < int(*rows)-2 ; j++{
			tablero[i][j] = CELULA_VIVA;
			fmt.Printf("insertando en: (%[1]d, %[2]d)\n",i,j);
		}	
	}*/

}


func semillaInicial(semilla int, columnas int, filas int, Tablero [][]int, celulas [][]Celula){
	rand.Seed(time.Now().UnixNano());
	for i := 0; i<semilla ; i++{
		col := rand.Intn(columnas - 0 + 1) + 0;
		fil := rand.Intn(filas - 0 + 1) + 0;
		celulas[(col/2)][(fil/2)+4].estado = CELULA_VIVA;
		Tablero[(col/2)][(fil/2)+4] = celulas[(col/2)][(fil/2)+4].estado;
		fmt.Printf("insertando en: (%[1]d, %[2]d)\n",col,fil);	
	}
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