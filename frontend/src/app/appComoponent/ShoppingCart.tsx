"use client";

import React, { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import ProductCard from "./ProductCard";
import { IProduct } from "./ProductCard.type";
import { useCart } from "./CartContext";

// // mock data
// const cart: IProduct[] = [
//   {
//     id: "1",
//     name: "Jordan Air Globe T-Shirt Kids",
//     price: 587000,
//     quantity: 1,
//     image:
//       "https://i1.t4s.cz/products/95d121-001/jordan-air-globe-t-shirt-kids-749837-95d121-001.png",
//     size: "XS",
//   },
//   {
//     id: "2",
//     name: "adidas Basketball Select Tee White",
//     price: 789000,
//     quantity: 2,
//     image:
//       "https://www.cosmossport.gr/2869439-product_medium/adidas-basketball-select-tee.jpg",
//     size: "L",
//   },
// ];

const ShoppingCart = () => {
  const { cart, removeFromCart } = useCart();
  const [selectItems, setSelectItems] = useState<IProduct[]>([]);
  const [totalPrice, setTotalPrice] = useState<number>(0);
  const [isSelectAll, setIsSelectAll] = useState<boolean>(false);
  const [checkedItems, setCheckedItems] = useState<boolean[]>(
    new Array(cart.length).fill(false)
  );

  const handleSelectAll = () => {
    const newCheckedItems = new Array(cart.length).fill(!isSelectAll);
    setCheckedItems(newCheckedItems);
    setIsSelectAll(!isSelectAll);

    if (!isSelectAll) {
      setSelectItems(cart);
    } else {
      setSelectItems([]);
    }
  };

  const handleCheckboxChange = (index: number) => {
    const newCheckedItems = [...checkedItems];
    newCheckedItems[index] = !newCheckedItems[index];
    setCheckedItems(newCheckedItems);

    const selected = cart.filter((_, idx) => newCheckedItems[idx]);
    setSelectItems(selected);
    setIsSelectAll(selected.length === cart.length);
  };

  const getTotalBill = useCallback((): number => {
    return selectItems.reduce(
      (total, item) => total + item.price * item.quantity,
      0
    );
  }, [selectItems]);

  useEffect(() => {
    setTotalPrice(getTotalBill());
  }, [selectItems, getTotalBill, cart]);

  return (
    <>
      <div className="flex flex-col gap-12">
        <div className="flex flex-col gap-6 h-max">
          {cart.length === 0 ? (
            <div className="flex justify-center">
              <span className="text-2xl">Shopping cart is empty. </span>
              <Link href="/">
                <span className="text-2xl text-blue-400 underline cursor-pointer">
                  Go shopping now!
                </span>
              </Link>
            </div>
          ) : (
            cart.map((product, index) => (
              <div
                key={product.id}
                className="flex flex-row items-center gap-20 px-16 py-6 bg-[#EBEBD5] h-full"
              >
                <div className="size-fit">
                  <input
                    type="checkbox"
                    checked={checkedItems[index]}
                    onChange={() => handleCheckboxChange(index)}
                    className="size-6 bg-white"
                  />
                </div>
                <ProductCard key={product.id} product={product} />
                <button
                  className="px-4 py-3 bg-white text-[#C6393F] self-end min-w-fit"
                  onClick={() => removeFromCart(product.id)}
                >
                  <span>Delete</span>
                </button>
              </div>
            ))
          )}
        </div>
        <div className="flex px-16 py-6 bg-[#EBEBD5] w-full">
          <div className="flex items-center justify-between w-full">
            <div className="flex items-center gap-4">
              <input
                type="checkbox"
                checked={isSelectAll}
                onChange={handleSelectAll}
                className="size-6 bg-white"
              />
              <label className="text-lg">Choose all</label>
            </div>
            <div>
              <p>
                Total bill:{" "}
                <span className="text-xl">
                  {Intl.NumberFormat("vi-VN").format(totalPrice)}{" "}
                  <span className="underline">đ</span>
                </span>
              </p>
            </div>
            <div>
              <Link href="/checkout">
                <button className="px-6 py-3 bg-[#EFD471] text-[#1A3C73]">
                  <span>Checkout</span>
                </button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default ShoppingCart;
