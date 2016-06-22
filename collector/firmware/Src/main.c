/**
  ******************************************************************************
  * File Name          : main.c
  * Description        : Main program body
  ******************************************************************************
  *
  * COPYRIGHT(c) 2016 STMicroelectronics
  *
  * Redistribution and use in source and binary forms, with or without modification,
  * are permitted provided that the following conditions are met:
  *   1. Redistributions of source code must retain the above copyright notice,
  *      this list of conditions and the following disclaimer.
  *   2. Redistributions in binary form must reproduce the above copyright notice,
  *      this list of conditions and the following disclaimer in the documentation
  *      and/or other materials provided with the distribution.
  *   3. Neither the name of STMicroelectronics nor the names of its contributors
  *      may be used to endorse or promote products derived from this software
  *      without specific prior written permission.
  *
  * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
  * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
  * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
  * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
  * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
  * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
  * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
  *
  ******************************************************************************
  */
/* Includes ------------------------------------------------------------------*/
#include "stm32f1xx_hal.h"

/* USER CODE BEGIN Includes */


#include <string.h>
#include "EtherShield.h"

#define SENSORS		16
#define MAX_CHANNELS	16

#define NET_BUF_SIZE (1<<10)
#define UART_BUF_SIZE (NET_BUF_SIZE << 2)

#define MAX(a,b) ((a)>(b)?(a):(b))
#define BUF_SIZE MAX(NET_BUF_SIZE, UART_BUF_SIZE)

typedef uint16_t dataitem_t;

struct sensorcommand {
	dataitem_t	command_id;
	dataitem_t	 sensor_id;
	dataitem_t	channels;
	dataitem_t	channel[MAX_CHANNELS];
};
typedef struct sensorcommand sensorcommand_t;

struct collectorcommand {
	dataitem_t	command_id;
	dataitem_t	 sensor_id;
};
typedef struct collectorcommand collectorcommand_t;

enum command_id {
	CMD_SETDATA,
	CMD_GETDATA,
};

/* USER CODE END Includes */

/* Private variables ---------------------------------------------------------*/
SPI_HandleTypeDef hspi1;
DMA_HandleTypeDef hdma_spi1_rx;
DMA_HandleTypeDef hdma_spi1_tx;

UART_HandleTypeDef huart1;
DMA_HandleTypeDef hdma_usart1_rx;
DMA_HandleTypeDef hdma_usart1_tx;

/* USER CODE BEGIN PV */
/* Private variables ---------------------------------------------------------*/

/* USER CODE END PV */

/* Private function prototypes -----------------------------------------------*/
void SystemClock_Config(void);
static void MX_GPIO_Init(void);
static void MX_DMA_Init(void);
static void MX_SPI1_Init(void);
static void MX_USART1_UART_Init(void);

/* USER CODE BEGIN PFP */
/* Private function prototypes -----------------------------------------------*/

void error (float error_num, char infinite) {
	//printf("%u\r\n", error_num);
	if (infinite)
		while (1) {
			int i = 0;
			while (i++ < ((int)((float)(error_num) / 1) + 1) ) {
				GPIOC->BSRR = LED_R_Pin;
				HAL_Delay(1000 / error_num);
				GPIOC->BSRR = LED_R_Pin << 16;
				HAL_Delay(1000 / error_num);
			}
		};

	GPIOC->BSRR = LED_R_Pin;
	HAL_Delay(1000);
	GPIOC->BSRR = LED_R_Pin << 16;
	HAL_Delay(100);

	int i=0;
	while (i++ < error_num) {
		GPIOC->BSRR = LED_R_Pin;
		HAL_Delay(300);
		GPIOC->BSRR = LED_R_Pin << 16;
		HAL_Delay(200);
	}

	HAL_Delay(500);
	NVIC_SystemReset();

	return;
}

static inline void blink(int times, int delay) {
	int i = 0;
	while (i++ < times) {
		GPIOC->BSRR = LED_R_Pin;
		HAL_Delay(delay);
		GPIOC->BSRR = LED_R_Pin << 16;
		HAL_Delay(delay);
	}

	return;
}

void ES_PingCallback(void) {
}

uint16_t get_udp_data_len(uint8_t *buf)
{
	int16_t i;
	i=(((int16_t)buf[IP_TOTLEN_H_P])<<8)|(buf[IP_TOTLEN_L_P]&0xff);
	i-=IP_HEADER_LEN;
	i-=8;
	if (i<=0){
		i=0;
	}
	return((uint16_t)i);
}

static uint16_t info_data_len = 0;
uint16_t packetloop_icmp_udp(uint8_t *buf,uint16_t plen)
{
	if(eth_type_is_arp_and_my_ip(buf,plen)){
		if (buf[ETH_ARP_OPCODE_L_P]==ETH_ARP_OPCODE_REQ_L_V){
			// is it an arp request 
			make_arp_answer_from_request(buf);
		}
		return(0);

	}
	// check if ip packets are for us:
	if(eth_type_is_ip_and_my_ip(buf,plen)==0){
		return(0);
	}

	if(buf[IP_PROTO_P]==IP_PROTO_ICMP_V && buf[ICMP_TYPE_P]==ICMP_TYPE_ECHOREQUEST_V){
		make_echo_reply_from_request(buf,plen);
		return(0);
	}

	if (buf[IP_PROTO_P]==IP_PROTO_UDP_V) {
		info_data_len=get_udp_data_len(buf);
		return(IP_HEADER_LEN+8+14);
	}

	return(0);
}


/* USER CODE END PFP */

/* USER CODE BEGIN 0 */

#define NET_HEADERS_LENGTH (ETH_HEADER_LEN + IP_HEADER_LEN + UDP_HEADER_LEN)
#define RECV_TIMEOUT 1000

static uint8_t ticked = 0;
uint8_t net_buf[NET_BUF_SIZE + 1];

void main_tick() {
	ticked = 1;

	return;
}

void HAL_Delay(__IO uint32_t Delay)
{
	uint32_t tickstart = 0;
	tickstart = HAL_GetTick();
	while((HAL_GetTick() - tickstart) < Delay)
	{
		packetloop_icmp_udp(net_buf, ES_enc28j60PacketReceive(NET_BUF_SIZE, net_buf));
	}
}

/* USER CODE END 0 */

int main(void)
{

  /* USER CODE BEGIN 1 */

  /* USER CODE END 1 */

  /* MCU Configuration----------------------------------------------------------*/

  /* Reset of all peripherals, Initializes the Flash interface and the Systick. */
  HAL_Init();

  /* Configure the system clock */
  SystemClock_Config();

  /* Initialize all configured peripherals */
  MX_GPIO_Init();
  MX_DMA_Init();
  MX_SPI1_Init();
  MX_USART1_UART_Init();

  /* USER CODE BEGIN 2 */
	uint8_t  uart_recvbuf[UART_BUF_SIZE], uart_sendbuf[UART_BUF_SIZE];
	sensorcommand_t    *scmd =    (sensorcommand_t *)uart_recvbuf;
	collectorcommand_t *ccmd = (collectorcommand_t *)uart_sendbuf;

	{
		HAL_StatusTypeDef r = HAL_UART_Receive_DMA(&huart1, uart_recvbuf, UART_BUF_SIZE);
		if (r != HAL_OK)
			error(4+r, 0);
	}

	uint8_t   local_mac[] = {0x02, 0x03, 0x04, 0x05, 0x06, 0x08};
	uint8_t  remote_mac[] = {0x00, 0x1b, 0x21, 0x39, 0x37, 0x26};
	uint8_t   local_ip[]  = {10,  4, 33, 124};
	uint8_t  remote_ip[]  = {10,  4, 33, 242};

	ES_enc28j60SpiInit(&hspi1);
	ES_enc28j60Init(local_mac);

	uint8_t enc28j60_rev = ES_enc28j60Revision();
	if (enc28j60_rev <= 0)
		error(2, 0);

	ES_init_ip_arp_udp_tcp(local_mac, local_ip, 80);

  /* USER CODE END 2 */

  /* Infinite loop */
  /* USER CODE BEGIN WHILE */
	/* // ICMP echo server
	while (1)
	{
		packetloop_icmp_udp(net_buf, ES_enc28j60PacketReceive(NET_BUF_SIZE, net_buf));
	}*/


	while (1)
	{
		uint8_t awaitingForSending = 0;
		dataitem_t channels;
		dataitem_t *net_data = (dataitem_t *)(&net_buf[NET_HEADERS_LENGTH]);

		if (awaitingForSending) {
			ES_send_udp_data2(net_buf, remote_mac, NET_HEADERS_LENGTH + channels*sizeof(*net_data), 26524, remote_ip, 36400);
			awaitingForSending = 0;
		}
		packetloop_icmp_udp(net_buf, ES_enc28j60PacketReceive(NET_BUF_SIZE, net_buf));

		{
			static uint8_t  awaitingForReceive = 0;
			static uint16_t timeoutCounter;
			static dataitem_t sensor_id = ~0;

			if (awaitingForReceive) {
				timeoutCounter++;
				if (timeoutCounter > RECV_TIMEOUT) {
					//error(3, 0);
					blink(3, 100);
					__HAL_UART_FLUSH_DRREGISTER(&huart1);
					awaitingForReceive = 0;
					GPIOC->BSRR = LED_G_Pin << 16;
					continue;
				}

				uint16_t itemsReceived = (UART_BUF_SIZE - huart1.hdmarx->Instance->CNDTR) / sizeof(dataitem_t);
				if (itemsReceived > 0) {
					if (scmd->command_id != CMD_SETDATA) {
						error(3, 0);
					}
				}
				if (itemsReceived > 1) {
					if (sensor_id != scmd->sensor_id) {
						error(4, 0);
					}
				}
				if (itemsReceived > 2) {
					channels  = scmd->channels;
				}

				if (itemsReceived-1 == channels) {
					net_data[0] = sensor_id;
					net_data[1] = channels;
					int i = 0;
					while (i < channels) {
						net_data[i+2] = scmd->channel[i];
						i++;
					}
					__HAL_UART_FLUSH_DRREGISTER(&huart1);
					awaitingForReceive = 0;
					awaitingForSending = 1;
					GPIOC->BSRR = LED_G_Pin << 16;
				}

				__HAL_UART_FLUSH_DRREGISTER(&huart1);

				continue;
			}

			if (ticked) {
				ticked = 0;
				HAL_Delay(1);
				GPIOC->BSRR = LED_G_Pin;
				ccmd->command_id = CMD_GETDATA;
				if (sensor_id == SENSORS-1)
					sensor_id = 0;
				ccmd->sensor_id  = ++sensor_id;
				{
					HAL_StatusTypeDef r;
					r = HAL_UART_Transmit(&huart1, (uint8_t *)ccmd, sizeof(*ccmd), 0xff);
					if (r != HAL_OK)
						error(4+r, 0);
				}

				timeoutCounter     = 0;
				awaitingForReceive = 1;
				continue;
			}
		}
	}

  /* USER CODE END WHILE */

  /* USER CODE BEGIN 3 */
  /* USER CODE END 3 */

}

/** System Clock Configuration
*/
void SystemClock_Config(void)
{

  RCC_OscInitTypeDef RCC_OscInitStruct;
  RCC_ClkInitTypeDef RCC_ClkInitStruct;

  RCC_OscInitStruct.OscillatorType = RCC_OSCILLATORTYPE_HSE;
  RCC_OscInitStruct.HSEState = RCC_HSE_ON;
  RCC_OscInitStruct.HSEPredivValue = RCC_HSE_PREDIV_DIV1;
  RCC_OscInitStruct.PLL.PLLState = RCC_PLL_ON;
  RCC_OscInitStruct.PLL.PLLSource = RCC_PLLSOURCE_HSE;
  RCC_OscInitStruct.PLL.PLLMUL = RCC_PLL_MUL9;
  HAL_RCC_OscConfig(&RCC_OscInitStruct);

  RCC_ClkInitStruct.ClockType = RCC_CLOCKTYPE_SYSCLK|RCC_CLOCKTYPE_PCLK1;
  RCC_ClkInitStruct.SYSCLKSource = RCC_SYSCLKSOURCE_PLLCLK;
  RCC_ClkInitStruct.AHBCLKDivider = RCC_SYSCLK_DIV1;
  RCC_ClkInitStruct.APB1CLKDivider = RCC_HCLK_DIV2;
  RCC_ClkInitStruct.APB2CLKDivider = RCC_HCLK_DIV1;
  HAL_RCC_ClockConfig(&RCC_ClkInitStruct, FLASH_LATENCY_2);

  HAL_SYSTICK_Config(HAL_RCC_GetHCLKFreq()/1000);

  HAL_SYSTICK_CLKSourceConfig(SYSTICK_CLKSOURCE_HCLK);

  /* SysTick_IRQn interrupt configuration */
  HAL_NVIC_SetPriority(SysTick_IRQn, 0, 0);
}

/* SPI1 init function */
void MX_SPI1_Init(void)
{

  hspi1.Instance = SPI1;
  hspi1.Init.Mode = SPI_MODE_MASTER;
  hspi1.Init.Direction = SPI_DIRECTION_2LINES;
  hspi1.Init.DataSize = SPI_DATASIZE_8BIT;
  hspi1.Init.CLKPolarity = SPI_POLARITY_LOW;
  hspi1.Init.CLKPhase = SPI_PHASE_1EDGE;
  hspi1.Init.NSS = SPI_NSS_SOFT;
  hspi1.Init.BaudRatePrescaler = SPI_BAUDRATEPRESCALER_4;
  hspi1.Init.FirstBit = SPI_FIRSTBIT_MSB;
  hspi1.Init.TIMode = SPI_TIMODE_DISABLED;
  hspi1.Init.CRCCalculation = SPI_CRCCALCULATION_DISABLED;
  hspi1.Init.CRCPolynomial = 10;
  HAL_SPI_Init(&hspi1);

}

/* USART1 init function */
void MX_USART1_UART_Init(void)
{

  huart1.Instance = USART1;
  huart1.Init.BaudRate = 115200;
  huart1.Init.WordLength = UART_WORDLENGTH_8B;
  huart1.Init.StopBits = UART_STOPBITS_1;
  huart1.Init.Parity = UART_PARITY_NONE;
  huart1.Init.Mode = UART_MODE_TX_RX;
  huart1.Init.HwFlowCtl = UART_HWCONTROL_NONE;
  huart1.Init.OverSampling = UART_OVERSAMPLING_16;
  HAL_UART_Init(&huart1);

}

/** 
  * Enable DMA controller clock
  */
void MX_DMA_Init(void) 
{
  /* DMA controller clock enable */
  __HAL_RCC_DMA1_CLK_ENABLE();

  /* DMA interrupt init */
  HAL_NVIC_SetPriority(DMA1_Channel2_IRQn, 0, 0);
  HAL_NVIC_EnableIRQ(DMA1_Channel2_IRQn);
  HAL_NVIC_SetPriority(DMA1_Channel3_IRQn, 0, 0);
  HAL_NVIC_EnableIRQ(DMA1_Channel3_IRQn);
  HAL_NVIC_SetPriority(DMA1_Channel4_IRQn, 0, 0);
  HAL_NVIC_EnableIRQ(DMA1_Channel4_IRQn);
  HAL_NVIC_SetPriority(DMA1_Channel5_IRQn, 0, 0);
  HAL_NVIC_EnableIRQ(DMA1_Channel5_IRQn);

}

/** Configure pins as 
        * Analog 
        * Input 
        * Output
        * EVENT_OUT
        * EXTI
*/
void MX_GPIO_Init(void)
{

  GPIO_InitTypeDef GPIO_InitStruct;

  /* GPIO Ports Clock Enable */
  __GPIOC_CLK_ENABLE();
  __GPIOD_CLK_ENABLE();
  __GPIOA_CLK_ENABLE();

  /*Configure GPIO pins : LED_R_Pin LED_G_Pin LED_B_Pin */
  GPIO_InitStruct.Pin = LED_R_Pin|LED_G_Pin|LED_B_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_OUTPUT_PP;
  GPIO_InitStruct.Speed = GPIO_SPEED_LOW;
  HAL_GPIO_Init(GPIOC, &GPIO_InitStruct);

  /*Configure GPIO pin : SPI1_CS_Pin */
  GPIO_InitStruct.Pin = SPI1_CS_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_OUTPUT_PP;
  GPIO_InitStruct.Speed = GPIO_SPEED_LOW;
  HAL_GPIO_Init(SPI1_CS_GPIO_Port, &GPIO_InitStruct);

}

/* USER CODE BEGIN 4 */

/* USER CODE END 4 */

#ifdef USE_FULL_ASSERT

/**
   * @brief Reports the name of the source file and the source line number
   * where the assert_param error has occurred.
   * @param file: pointer to the source file name
   * @param line: assert_param error line source number
   * @retval None
   */
void assert_failed(uint8_t* file, uint32_t line)
{
  /* USER CODE BEGIN 6 */
  /* User can add his own implementation to report the file name and line number,
    ex: printf("Wrong parameters value: file %s on line %d\r\n", file, line) */
  /* USER CODE END 6 */

}

#endif

/**
  * @}
  */ 

/**
  * @}
*/ 

/************************ (C) COPYRIGHT STMicroelectronics *****END OF FILE****/
