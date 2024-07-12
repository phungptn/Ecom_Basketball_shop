import React from "react";
import Link from "next/link";
import Image from "next/image";
import { FetchProduct, Product } from "../../../lib/product";

<<<<<<< Updated upstream
export default function Card () {
    return (
        <>
        {products.map(product => (
          <div key={product.productID} className="bg-primary flex flex-col h-[400px] w-72">
            <Image 
              src={product.imageURL} 
              alt={product.productName} 
              width={150} 
              height={150} 
              className="p-5 min-w-72 object-cover"
              loading="lazy" 
            />
            <span className="text-base-content font-semibold w-72 px-5 text-lg">
              {product.productName}
            </span>
            <div className="flex flex-row items-center justify-between w-72 px-5 mb-5">
              <span className="text-base-content text-lg">
                ${product.price}
              </span>
              <button className="btn bg-secondary rounded-none text-neutral font-semibold text-xl min-w-32 hover:bg-accent hover:text-primary">
                Xem
              </button>
            </div>
          </div>
        ))}
      </>
    )
}
=======
interface CardProps {
  limit?: number;
}

export default async function Card({ limit }: CardProps) {
  let products: Array<Product> = [];

  try {
    products = await FetchProduct();
    console.log("product:", products);
  } catch (error) {
    console.error('Error fetching products:', error);
  }

  // Apply the limit if provided
  const displayedProducts = limit ? products.slice(0, limit) : products;

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN').format(price);
  };

  return (
    <>
      {displayedProducts.map(product => (
        <div key={product.productID} className="bg-primary flex flex-col h-64 w-40 sm:h-80 sm:w-56 lg:h-[410px] lg:w-72">
          <Image 
            src={product.imageURL} 
            alt={product.productName} 
            width={150} 
            height={150} 
            className="p-5 w-48 sm:w-56 lg:w-72 object-cover"
            loading="lazy" 
          />
          <span className="text-base-content font-semibold w-40 sm:w-56 lg:w-72 h-14 px-5 text-xs sm:text-sm lg:text-lg">
            {product.productName}
          </span>
          <div className="flex flex-row items-center justify-between w-40 sm:w-56 lg:w-72 px-5 pb-1">
            <span className="text-base-content text-xs sm:text-sm lg:text-lg">
              đ {formatPrice(product.price)}
            </span>
            <button className="btn bg-secondary rounded-none text-neutral font-semibold text-xs sm:text-sm lg:text-lg w-12 sm:w-24 lg:w-32 sm:h-4 lg:h-12 hover:bg-accent hover:text-primary">
              Xem
            </button>
          </div>
        </div>
      ))}
    </>
  );
}
>>>>>>> Stashed changes
